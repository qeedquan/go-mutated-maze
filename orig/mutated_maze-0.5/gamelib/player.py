__author__ = 'sonneveld'

import sys
import os
import random
import pygame
import pygame.image

import ega
from mob import Mob
import data

from surfutil import  load_sprite_surface


FOLLOW_DELAY = 6
MAX_FOLLOWS = 16


class RailsThing(object):


    def __init__(self, gamestate, pos):
        '''
        :type for_maze: Maze
        '''

        self.gamestate = gamestate
        self.for_maze = gamestate.maze
        self.surface = load_sprite_surface("hero_dude.png", ega.SPRITE_COLOURS)

        self.death_snd = pygame.mixer.Sound(data.filepath("death.wav"))
        self.pickup_friend_snd = pygame.mixer.Sound(data.filepath("pickup_friend.wav"))

        self.x,self.y = pos
        self.vx = 0
        self.vy = 0

        self.wanted_direction=None
        self.direction=None

        # fill up history with the same pos
        self.pos_history = [ (self.x, self.y) ] * FOLLOW_DELAY*MAX_FOLLOWS

        self.followers = []

        self.is_dead = False

    @property
    def hitbox(self):
        return pygame.Rect(self.x+4, self.y+3, 8, 10)

    def on_mob_hit(self, mob):
        if mob.mob_type is 'friendly':
            if mob not in self.followers:
                self.pickup_friend_snd.play()
                mob.follow(self)
        elif mob.mob_type is 'enemy':
            self.death_snd.play()
            self.die(mob)

    def die(self, killing_mob):
        # try to kill one of our followers first
        if killing_mob:
            for x in self.followers:
                if type(x) is Mob:
                    self.unfollow(x)
                    x.die(killing_mob)
                    killing_mob.die(killing_mob)
                    return

        self.is_dead = True
        self.death_vx = random.choice( (-3, -2, -1, 1, 2, 3))
        self.death_vy = -3

    @property
    def mob_followers(self):
        return [x for x in self.followers if type(x) is Mob]

    def add_pos_history(self, pos):

        last = self.pos_history[-1]
        if pos == last:
            return
        self.pos_history.append( pos )
        self.pos_history = self.pos_history[-FOLLOW_DELAY*MAX_FOLLOWS:]


    def register_follower(self, f):
        if f in self.followers:
            return
        self.followers.append(f)

    def unfollow(self, f):
        if f not in self.followers:
            return
        self.followers.remove(f)

    def get_follower_index(self, f):
        return self.followers.index(f)

    def get_follower_pos(self, f):
        hist_index = FOLLOW_DELAY*(self.get_follower_index(f)+1)
        pos = self.pos_history[-hist_index]
        return pos


    def update_dead(self):
        self.surface = pygame.transform.rotate(self.surface, 90)
        self.x += self.death_vx
        self.y += self.death_vy
        self.death_vy += 1

    def update_alive(self):

        pressed = pygame.key.get_pressed()
        if pressed[pygame.K_UP]:
            self.wanted_direction="up"
        if pressed[pygame.K_DOWN]:
            self.wanted_direction="down"
        if pressed[pygame.K_LEFT]:
            self.wanted_direction="left"
        if pressed[pygame.K_RIGHT]:
            self.wanted_direction="right"


        # allow change in opposite direction
        if  (self.direction is 'up' and self.wanted_direction is 'down') or \
            (self.direction is 'down' and self.wanted_direction is 'up') or \
            (self.direction is 'left' and self.wanted_direction is 'right') or \
            (self.direction is 'right' and self.wanted_direction is 'left'):
            self.direction = self.wanted_direction

        # otherwise only change if right in middle of node
        if self.x % 16 == 0 and self.y % 16 == 0:
            n = self.for_maze.get_node(self.x//16, self.y//16)
            if  (self.wanted_direction is 'up' and n.up_open) or \
                (self.wanted_direction is 'down' and n.down_open) or \
                (self.wanted_direction is 'left' and n.left_open) or \
                (self.wanted_direction is 'right' and n.right_open):
                self.direction = self.wanted_direction

        dx,dy=0,0
        if self.direction is 'up':
            dy = -1
            assert self.x%16 == 0
            if self.y%16 == 0:
                maze_x = self.x//16
                maze_y = self.y//16
                n = self.for_maze.get_node(maze_x, maze_y)
                if not n.up_open:
                    dy=0
        elif self.direction is 'down':
            dy = 1
            assert self.x%16 == 0
            if self.y%16 == 0:
                maze_x = self.x//16
                maze_y = self.y//16
                n = self.for_maze.get_node(maze_x, maze_y)
                if not n.down_open:
                    dy=0
        elif self.direction is 'left':
            dx = -1
            assert self.y%16 == 0
            if self.x%16 == 0:
                maze_x = self.x//16
                maze_y = self.y//16
                n = self.for_maze.get_node(maze_x, maze_y)
                if not n.left_open:
                    dx=0
        elif self.direction is 'right':
            dx = 1
            assert self.y%16 == 0
            if self.x%16 == 0:
                maze_x = self.x//16
                maze_y = self.y//16
                n = self.for_maze.get_node(maze_x, maze_y)
                if not n.right_open:
                    dx=0

        self.x+= dx * 2
        self.y+= dy * 2



    def update(self):

        """
        the aim here is for 'pacman' like logic where you keep on going but
        you can suggest where you'd like to go if there is an opening.
        """

        if self.is_dead:
            self.update_dead()
        else:
            self.update_alive()

        self.add_pos_history( (self.x, self.y))




    def render(self, to_surface):
        to_surface.blit(self.surface, (self.x, self.y))




