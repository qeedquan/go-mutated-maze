'''Game main module.

Contains the entry point used by the run_game.py script.

Feel free to put all your game code here, or in other modules in this "gamelib"
package.
'''

import sys
import os
import random
import pygame
import pygame.image
from gamelib import surfutil

import maze
import ega
import data
from player import RailsThing
from mob import Mob
from fog import Fog
from textutil import TextBufferRenderer

DISPLAY_SIZE = (480,320)
BG_COLOUR = 0x1c1c1c
SCORE_POS = (380, 25)



# node distance from player for spawning mobs
INIT_DISTANCE_FROM_PLAYER = 5

# number of mobs to spawn
INIT_NUM_OF_MOBS = 15

GOOD_MOBS = 4
BAD_MOBS = 4

MAX_CHARGE = 10000
CHARGE_PER_FRAME = 50

SPARK_TIME_MS = 300


class Key(object):

    def __init__(self, gamestate, pos):
        '''
        :type for_maze: Maze
        '''

        self.gamestate = gamestate
        self.for_maze = gamestate.maze
        p = data.filepath("key.png")
        i = pygame.image.load(p)
        i = i.convert_alpha()
        self.surface = i

        self.x,self.y = pos

        self.follow_sprite = None

    @property
    def hitbox(self):
        return pygame.Rect(self.x+4, self.y+4, 8, 8)
    
    def follow(self, other_sprite):
        self.follow_sprite = other_sprite
        self.follow_index = other_sprite.register_follower(self)

    def update(self):
        if self.follow_sprite is None:
            return
        self.x,self.y = self.follow_sprite.get_follower_pos(self)

    def render(self, to_surface):
        to_surface.blit(self.surface, (self.x,self.y))



class LockedDoor(object):

    def __init__(self, gamestate, pos):
        '''
        :type for_maze: Maze
        '''

        self.gamestate = gamestate
        self.for_maze = gamestate.maze
        p = data.filepath("lock.png")
        i = pygame.image.load(p)
        i = i.convert_alpha()
        self.surface = i

        self.x, self.y = pos


    def update(self):
        pass

    def render(self, to_surface):
        to_surface.blit(self.surface, (self.x,self.y))


class GameState(object):

    def __init__(self):
        self.reset_game()
        self.sparks = []

    @property
    def all(self):
        if self.player:
            yield self.player
        for x in self.mobs:
            yield x

    @property
    def alive_mobs(self):
        return [x for x in self.mobs if not x.is_dead]

    def add_spark(self, pos):
        s = (pygame.time.get_ticks(), pos)
        self.sparks.append(s)

    def reset_game(self):
        self.maze = None
        self.mobs = []
        self.sparks = []
        self.player = None

        self.score = 0
        self.time = 60
        self.time_left = None
        self.level = 0

        self.mobs_saved = 0
        self.mobs_available = 0

        self.mobs_saved_total = 0


    def next_level(self):
        self.maze = None
        self.mobs = []
        self.player = None
        self.level += 1
        self.time = max(45, 120 - (self.level-1)*20)
        self.mobs_available = INIT_NUM_OF_MOBS + (self.level-1)*3
        self.time_left = None
        self.charge = MAX_CHARGE



def title(display_surface, gamestate):
    display_surface.fill(ega.BLACK)
    pygame.display.flip()

    titlefont = pygame.font.SysFont("Arial", 30)
    titlesurf = titlefont.render("Godspeed You! Mutated Maze", True, ega.BRIGHT_MAGENTA)

    instrfont = pygame.font.SysFont("Arial", 15)
    namesurf = instrfont.render("Copyright 2011 Nick Sonneveld.", True, ega.BRIGHT_MAGENTA)
    instrsurf = instrfont.render("Press <SPACE> to start.", True, ega.BRIGHT_MAGENTA)

    clock = pygame.time.Clock()
    while True:

        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return 'quit'
            elif event.type == pygame.KEYDOWN and event.key == pygame.K_ESCAPE:
                return 'quit'
            elif event.type == pygame.KEYDOWN and event.key == pygame.K_SPACE:
                gamestate.reset_game()
                return 'instructions'

        display_surface.fill(ega.BLACK)

        title_x = (display_surface.get_width() - titlesurf.get_width())/2
        display_surface.blit(titlesurf, (title_x,50))

        title_x = (display_surface.get_width() - namesurf.get_width())/2
        display_surface.blit(namesurf, (title_x,240))

        title_x = (display_surface.get_width() - instrsurf.get_width())/2
        display_surface.blit(instrsurf, (title_x,270))

        pygame.display.flip()
        clock.tick(30)



def youwin(display_surface, gamestate):
    display_surface.fill(ega.BLACK)
    pygame.display.flip()

    headerfont = pygame.font.SysFont("Arial", 40)
    textfont = pygame.font.SysFont("Arial", 15)


    titlesurf = headerfont.render("You ESCAPED!", True, ega.BRIGHT_MAGENTA)


    text_list = []

    # time left
    text_list.append(("Time Left: %d seconds"%gamestate.time_left, ega.BRIGHT_MAGENTA))
    score_time = gamestate.time_left * 30
    text_list.append(("+%d"%score_time, ega.BRIGHT_GREEN))
    text_list.append(("",ega.BRIGHT_MAGENTA))

    # people saved
    text_list.append(("Mobs saved: %d"%gamestate.mobs_saved, ega.BRIGHT_MAGENTA))
    score_mobs = gamestate.mobs_saved * 500
    text_list.append(("+%d"%score_mobs, ega.BRIGHT_GREEN))
    score_bonus = 0
    if gamestate.mobs_saved > 0 and gamestate.mobs_saved == gamestate.mobs_available:
        score_bonus = 2000
        text_list.append(("You saved everyone!", ega.BRIGHT_GREEN))
        text_list.append(("+%d"%score_bonus, ega.BRIGHT_GREEN))
        text_list.append(("You are a super player!", ega.BRIGHT_GREEN))
    text_list.append(("",ega.BRIGHT_MAGENTA))

    # = score
    gamestate.score += score_time + score_mobs + score_bonus
    text_list.append(("Score: %d"%gamestate.score, ega.BRIGHT_MAGENTA))
    text_list.append(("",ega.BRIGHT_MAGENTA))

    text_list.append(("Press <SPACE> to continue.",ega.BRIGHT_MAGENTA))

    clock = pygame.time.Clock()
    while True:

        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return 'quit'
            elif event.type == pygame.KEYDOWN and event.key == pygame.K_ESCAPE:
                return 'title'
            elif event.type == pygame.KEYDOWN and event.key == pygame.K_SPACE:
                return 'playgame'

        display_surface.fill(ega.BLACK)

        title_x = (display_surface.get_width() - titlesurf.get_width())/2
        display_surface.blit(titlesurf, (title_x,20))

        y = 100
        for text,colour in text_list:
            rendered_text = textfont.render(text, True, colour)
            x = (display_surface.get_width() - rendered_text.get_width())/2
            display_surface.blit(rendered_text, (x,y))
            y += rendered_text.get_height()

        pygame.display.flip()
        clock.tick(30)





def gameover_screen(display_surface, gamestate):
    titlefont = pygame.font.SysFont("Arial", 30)
    titlesurf = titlefont.render("Game over man, game over!", True, ega.BRIGHT_MAGENTA)

    textfont = pygame.font.SysFont("Arial", 15)

    text_list = []
    text_list.append(("Final score: %d"%gamestate.score, ega.BRIGHT_MAGENTA))
    text_list.append(("Mobs saved: %d"%gamestate.mobs_saved_total, ega.BRIGHT_MAGENTA))
    text_list.append(("",ega.BRIGHT_MAGENTA))
    text_list.append(("Press <SPACE> to continue.",ega.BRIGHT_MAGENTA))



    clock = pygame.time.Clock()
    while True:

        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return 'quit'
            elif event.type == pygame.KEYDOWN and event.key == pygame.K_ESCAPE:
                return 'title'
            elif event.type == pygame.KEYDOWN and event.key == pygame.K_SPACE:
                return 'title'

        display_surface.fill(ega.BLACK)

        title_x = (display_surface.get_width() - titlesurf.get_width())/2
        display_surface.blit(titlesurf, (title_x,50))

        y = 200
        for text,colour in text_list:
            rendered_text = textfont.render(text, True, colour)
            x = (display_surface.get_width() - rendered_text.get_width())/2
            display_surface.blit(rendered_text, (x,y))
            y += rendered_text.get_height()

        pygame.display.flip()
        clock.tick(30)



def instructions_screen(display_surface, gamestate):

    instructions_surface = surfutil.load_sprite_surface("instructions.png")

    clock = pygame.time.Clock()
    while True:

        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return 'quit'
            elif event.type == pygame.KEYDOWN and event.key == pygame.K_ESCAPE:
                return 'title'
            elif event.type == pygame.KEYDOWN and event.key == pygame.K_SPACE:
                return 'playgame'

        display_surface.blit(instructions_surface, (0,0))
        pygame.display.flip()
        clock.tick(30)



def playgame(screen, gamestate):

    gamestate.next_level()

    opendoor_snd = pygame.mixer.Sound(data.filepath("open_door.wav"))
    pickup_key_snd = pygame.mixer.Sound(data.filepath("pickup_key.wav"))

    background = pygame.Surface(screen.get_size())
    background = background.convert()
    background.fill(BG_COLOUR)

    timeupfont = pygame.font.SysFont("Arial", 30)
    
    tbuffer = TextBufferRenderer("Arial", 15)

    m = maze.Maze(22, 18)
    m.generate()
    maze_surface = pygame.Surface((m.width*16, m.height*16), pygame.constants.SRCALPHA, 32).convert_alpha()

    gamestate.maze = m

    #fixed positions
    player_start_node = m.get_node(0,0)
    player_sprite = RailsThing(gamestate, player_start_node.pos_px)
    gamestate.player = player_sprite
    f = Fog(gamestate)

    lock_start_node = m._data[-1]
    lockspr = LockedDoor(gamestate, lock_start_node.pos_px)

    # ensure everything else is on its own square
    # but not too close to player

    availnodes = [x for x in m._data if x.distance_from(player_start_node) > INIT_DISTANCE_FROM_PLAYER]
    random.shuffle(availnodes)

    for x in xrange(gamestate.mobs_available):
        n = availnodes.pop()
        mob = Mob(gamestate, n.pos_px)
        if random.choice([True, False, False]):
            mob.toggle_mob_type()
        gamestate.mobs.append(mob)


    availkeynodes = [x for x in m._data if \
                     x.distance_from(player_start_node) > INIT_DISTANCE_FROM_PLAYER and \
                     x.distance_from(lock_start_node) > INIT_DISTANCE_FROM_PLAYER ]
    key_node = random.choice(availkeynodes)
    keyspr = Key(gamestate, key_node.pos_px)


    mutator_surface = pygame.Surface((64, 64), pygame.constants.SRCALPHA, 32).convert_alpha()
    mutator_surface.fill((255, 255, 255, 30))

    charge_surface = surfutil.load_sprite_surface("chargebar.png")
    charge_surface = pygame.transform.scale(charge_surface, (charge_surface.get_width()*2, charge_surface.get_height()*2))

    spark_surface = surfutil.load_sprite_surface("spark.png")

    # blank screen
    screen.blit(background, (0, 0))
    pygame.display.flip()

    clock = pygame.time.Clock()
    start_time = pygame.time.get_ticks()
    MAZE_POS = (16, 16)
    death_count = 0
    raw_click_pos = (0,0)
    while True:

        # USER INPUT
        # ======================================

        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return 'quit'
            elif event.type == pygame.KEYDOWN and event.key == pygame.K_ESCAPE:
                return 'title'
            elif event.type == pygame.KEYDOWN and event.key == pygame.K_g:
                print 'No soup for you GruikInc!'
            #    m.generate()
            #    drawmaze(maze_surface, m)
            elif event.type == pygame.MOUSEBUTTONDOWN and event.button == 1:
                raw_click_pos = pygame.mouse.get_pos()
                #print raw_click_pos
                click_pos = (raw_click_pos[0]-MAZE_POS[0], raw_click_pos[1]-MAZE_POS[1])

                if gamestate.charge == MAX_CHARGE:
                    f.mutate_snd.play()
                    r = pygame.Rect(click_pos, mutator_surface.get_size())
                    selected = set(m.collides_nodes(r))
                    m.regenerate_selected(selected)
                    gamestate.charge = 0

        # UPDATE
        # ======================================

        current_time = pygame.time.get_ticks()
        gamestate.time_left = max(int(gamestate.time - (current_time-start_time)/1000), 0)
        if gamestate.time_left <= 0 and not player_sprite.is_dead:
            player_sprite.die(None)

        player_sprite.update()
        keyspr.update()
        lockspr.update()

        player_hitbox = player_sprite.hitbox
        for mob in gamestate.mobs:
            mob.update()

        for mob in gamestate.alive_mobs:
            if player_hitbox.colliderect(mob.hitbox):
                player_sprite.on_mob_hit(mob)

        if keyspr.follow_sprite is not player_sprite and player_hitbox.colliderect(keyspr.hitbox):
            pickup_key_snd.play()
            gamestate.score+=200
            keyspr.follow(player_sprite)

        if keyspr.follow_sprite is player_sprite and lockspr.x == player_sprite.x and lockspr.y == player_sprite.y:
            opendoor_snd.play()
            gamestate.score+=200
            gamestate.mobs_saved = len(player_sprite.mob_followers)
            gamestate.mobs_saved_total += gamestate.mobs_saved
            return 'win'

        f.update()
        if f.passed:
            f = Fog(gamestate)

        if player_sprite.is_dead:
            death_count += 1
            if death_count > 100:
                return 'gameover'

        gamestate.charge = min(MAX_CHARGE, gamestate.charge + CHARGE_PER_FRAME)

        current_time = pygame.time.get_ticks()
        gamestate.sparks = [x for x in gamestate.sparks if  (current_time- x[0])<SPARK_TIME_MS]

        # RENDER
        # ======================================

        screen.blit(background, (0, 0))

        m.render_to_surface(maze_surface)
        lockspr.render(maze_surface)

        for mob in gamestate.mobs:
            mob.render(maze_surface)
        keyspr.render(maze_surface)

        # make sure we can always see your dude
        player_sprite.render(maze_surface)

        # fog
        f.render(maze_surface)

        #sparks
        for (t,p) in gamestate.sparks:
            maze_surface.blit(spark_surface, p)

        screen.blit(maze_surface, MAZE_POS)

        screen.blit(mutator_surface, pygame.mouse.get_pos())

        tbuffer.clear()
        tbuffer.add("Level: %d"%gamestate.level)
        tbuffer.add("Score: %d"%gamestate.score)

        # print the status
        c = ega.BRIGHT_MAGENTA
        if gamestate.time_left < 15 and gamestate.time_left&1:
            c = ega.BRIGHT_WHITE
        tbuffer.add("Time: %d"%gamestate.time_left, c)
        tbuffer.render_to_surface(screen, SCORE_POS)

        # print time up msg if needed
        if gamestate.time_left <= 0:
            timeupsurf = timeupfont.render("TIME UP!", True, ega.BRIGHT_MAGENTA)
            screen.blit(timeupsurf, (128, 100))

        charge_width = int(charge_surface.get_width() * gamestate.charge/MAX_CHARGE)
        charge_area = pygame.Rect(0, 0, charge_width, 100)
        screen.blit(charge_surface, (380, 280), charge_area)

        pygame.display.flip()

        # TIMING
        # ======================================

        clock.tick(30)




def main():

    #print "Hello from your game's main()"
    #print data.load('sample.txt').read()
    pygame.init()
    screen = pygame.display.set_mode(DISPLAY_SIZE)
    pygame.display.set_caption('Godspeed You! Mutated Maze')
    pygame.mouse.set_visible(0)

    state = 'title'
    gamestate = GameState()

    while state != 'quit' and state is not None:
        if state is 'title':
            state = title(screen, gamestate)
        elif state is 'instructions':
            state = instructions_screen(screen, gamestate)
        elif state is 'playgame':
            state = playgame(screen, gamestate)
        elif state is 'win':
            state = youwin(screen, gamestate)
        elif state is 'gameover':
            state = gameover_screen(screen, gamestate)
        else:
            raise Exception("unknown state: "+state)





if __name__ == "__main__":
    main()
