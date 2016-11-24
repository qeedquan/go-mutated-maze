import sys
import os
import random
import pygame
import pygame.image

import ega
import data
from surfutil import load_sprite_surface


EXCLAMATION_COUNT_INIT = 16


# mob can be good or bad
class Mob(object):


    def __init__(self, gamestate, pos):
        '''
        :type for_maze: Maze
        '''

        self.gamestate = gamestate
        self.for_maze = gamestate.maze

        self.good_surface = load_sprite_surface("other_dude.png", ega.SPRITE_COLOURS)
        self.bad_surface = load_sprite_surface("ghost.png", ega.BAD_COLOURS)
        self.exclamation_surface = load_sprite_surface("exclamation.png")

        self.exclamation_snd = pygame.mixer.Sound(data.filepath("exclamation.wav"))

        # will set type and surface and state
        self.set_friendly_mob_type()

        self.x, self.y = pos
        self.vx = 0
        self.vy = 0

        self.wanted_direction=None
        self.direction=None

        self.is_dead = False

    @property
    def hitbox(self):
        return pygame.Rect(self.x+4, self.y+3, 8, 10)

    def die(self, killing_mob):
        if self.mob_type is 'enemy':
            self.gamestate.score += 200

        self.is_dead = True
        self.death_vx = random.choice( (-3, -2, -1, 1, 2, 3))
        self.death_vy = -3

    # MOB TYPES
    # ===============================================================

    def toggle_mob_type(self):
        if self.mob_type is 'friendly':
            self.set_enemy_mob_type()
        elif self.mob_type is 'enemy':
            self.set_friendly_mob_type()





    # FOG MUTATIONS
    # ===============================================================

    def do_fog_mutate(self, pos):
        #self.gamestate.add_spark((self.x, self.y))
        self.gamestate.add_spark(pos)

        if self.state is 'follow_player':
            self.follow_sprite.unfollow(self)

        self.init_nothing()
        self.x,self.y = pos
        self.toggle_mob_type()


    # UPDATE
    # ===============================================================

    def update(self):


        if self.is_dead:
            self.surface = pygame.transform.rotate(self.surface, 90)
            self.x += self.death_vx
            self.y += self.death_vy
            self.death_vy += 1
        elif self.mob_type is 'friendly':
            self.friendly_update()
        elif self.mob_type is 'enemy':
            self.enemy_update()


    # FRIENDLY BEHAVIOR
    # ===============================================================

    def set_friendly_mob_type(self):
        self.mob_type = 'friendly'
        self.surface = self.good_surface
        self.init_nothing()

    def friendly_update(self):
        if self.state is 'nothing':
            random.choice((self.init_wait, self.init_random))()

        if self.state is 'random':
            self.update_random()
        elif self.state is 'wait':
            self.update_wait()
        elif self.state is 'follow_player':
            self.update_follow_player()

    def init_nothing(self):
        self.state = 'nothing'

    def init_wait(self):
        self.state = 'wait'
        self.wait_count = random.randint(3,8)

    def update_wait(self):
        self.wait_count -= 1
        if self.wait_count <= 0:
            self.init_nothing()

    def init_random(self):
        self.state = 'random'

        # otherwise only change if right in middle of node
        n = self.for_maze.get_node(self.x//16, self.y//16)
        avail_dir = list(n.avail_directions)
        self.wanted_direction = random.choice(avail_dir)
        self.direction = self.wanted_direction

    def update_random(self):
        speed = 1

        dx,dy = {
            'up' : (0, -1),
            'down': (0, 1),
            'left': (-1, 0),
            'right': (1, 0),
            None: (0,0),
        } [self.direction]

        self.x+= dx * speed
        self.y+= dy * speed

        # otherwise only change if right in middle of node
        if self.x % 16 == 0 and self.y % 16 == 0:
            self.init_nothing()
            return


    # FOLLOWING PLAYER
    # ===============================================================


    def follow(self, other_sprite):
        self.follow_sprite = other_sprite
        self.follow_index = other_sprite.register_follower(self)
        self.state = 'follow_player'


    def update_follow_player(self):
        if self.follow_sprite is None:
            return
        self.x,self.y = self.follow_sprite.get_follower_pos(self)


    # ENEMY MOVEMENT
    # ===============================================================

    def set_enemy_mob_type(self):
        self.mob_type = 'enemy'
        self.surface = self.bad_surface
        self.init_nothing()

    def enemy_update(self):
        if self.state is 'nothing':
            random.choice((self.init_wait, self.init_random))()

        self.detect_player()

        if self.state is 'random':
            self.update_random()
        elif self.state is 'wait':
            self.update_wait()
        elif self.state is 'exclamation':
            self.update_exclamation()
        elif self.state is 'chase':
            self.update_chase()

    def init_exclamation(self, vx, vy):
        self.state = 'exclamation'
        self.exclamation_count = EXCLAMATION_COUNT_INIT
        self.x >>= 1 # ensure multiple of 2
        self.x <<= 1
        self.y >>= 1
        self.y <<= 1
        self.chase_vx = vx
        self.chase_vy = vy
        self.exclamation_snd.play()

    def update_exclamation(self):
        self.exclamation_count -=1
        if self.exclamation_count <= 0:
            self.init_chase()


    def init_chase(self):
        self.state = 'chase'

    def update_chase(self):
        # if at intersection
        if self.x%16==0 and self.y%16 ==0:
            node = self.for_maze.get_node(self.x/16, self.y/16)
            if self.chase_vx < 0:
                if not node.left_open:
                    self.init_nothing()
                    return
            if self.chase_vx > 0:
                if not node.right_open:
                    self.init_nothing()
                    return
            if self.chase_vy < 0:
                if not node.up_open:
                    self.init_nothing()
                    return
            if self.chase_vy > 0:
                  if not node.down_open:
                    self.init_nothing()
                    return

        self.x += self.chase_vx
        self.y += self.chase_vy



    def detect_player(self):
        if self.state is 'exclamation':
            return

        # if on a row
        if self.y % 16 == 0:
            # do we see the player to the left or right?
            # if yes, set state to chasee
            if self.find_player_on(lambda x: x.left_node if x.left_open else None):
                if self.state is 'chase' and self.chase_vx < 0:
                    return
                self.init_exclamation(-2, 0)
                return
            if self.find_player_on(lambda x: x.right_node if x.right_open else None):
                if self.state is 'chase' and self.chase_vx > 0:
                    return
                self.init_exclamation(2, 0)
                return


        # if on a column
        if self.x % 16 == 0:
            # do we see player above or below?
            if self.find_player_on(lambda x: x.up_node if x.up_open else None):
                if self.state is 'chase' and self.chase_vy < 0:
                    return
                self.init_exclamation(0, -2)
                return
            if self.find_player_on(lambda x: x.down_node if x.down_open else None):
                if self.state is 'chase' and self.chase_vy > 0:
                    return
                self.init_exclamation(0, 2)
                return


    def find_player_on(self, getnextnode):

        player_hitbox = self.gamestate.player.hitbox

        current_node = self.for_maze.get_node(self.x/16, self.y/16)
        while current_node is not None:
            if current_node.hitbox.colliderect(player_hitbox):
                return True
            current_node = getnextnode(current_node)

        return False



    # RENDERING
    # ===============================================================

    def render(self, to_surface):
        to_surface.blit(self.surface, (self.x, self.y))
        if self.state is 'exclamation':
            to_surface.blit(self.exclamation_surface, (self.x, self.y-8))


