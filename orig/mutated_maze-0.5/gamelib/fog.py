
import sys
import os
import random
import pygame
import pygame.image
import itertools

import data
from gamelib import ega


FOGFILES=(
    ('fog01.png', 'fog01_inverse.png'),
    ('fog02.png', 'fog02_inverse.png'),
    ('fog03.png', 'fog03_inverse.png'),
    ('fog04.png', 'fog04_inverse.png'),
    ('fog05.png', 'fog05_inverse.png'),
    ('fog06.png', 'fog06_inverse.png'),
)

# setup a random list of colours
maze_colours = list(ega.SPRITE_COLOURS)
random.shuffle(maze_colours)
colour_iterator = itertools.cycle(maze_colours)

class Fog(object):

    def __init__(self, gamestate):
        '''
        :type gamestate: GameState
        '''

        self.gamestate = gamestate
        self.for_maze = gamestate.maze

        (fog_name, inverse_name) = random.choice(FOGFILES)
        self.fog_surface = self.load_fog(fog_name)
        self.inverse_surface = self.load_fog(inverse_name)

        self.surface = self.fog_surface

        self.mutate_snd = pygame.mixer.Sound(data.filepath("mutate.wav"))

        self.x = gamestate.maze.width_px
        self.y = 0
        self.vx = -1
        self.vy = 0

        self.wanted_direction=None
        self.direction=None

        self.passed = False

        # figure this out before hand.
        #self.mutate_point = random.randint(-self.for_maze.width_px/2, self.for_maze.width_px/2)

        self.mutated = False

        self.mutate_count = random.randint(self.for_maze.width_px*0.5, self.for_maze.width_px*1.5)
        self.flash_count = 0



    def load_fog(self, path):
        maze = self.gamestate.maze
        datapath = data.filepath(path)
        i = pygame.image.load(datapath)
        i = pygame.transform.scale(i, (maze.width_px, maze.height_px))
        i = i.convert_alpha()
        return i

    def mutate(self):
        self.mutate_snd.play()

        maze = self.gamestate.maze

        selected_nodes = []

        new_maze_colour = colour_iterator.next()

        # select nodes under the fog
        for node in maze._data:
            c_x = node.x*16 + 8 # mid point
            c_y = node.y*16 + 8
            if self.pos_covered_with_fog((c_x, c_y)):
                node.colour = new_maze_colour
                selected_nodes.append(node)

        # find all the mobs and move them, also mutate them
        mutate_mobs = []
        for mob in self.gamestate.alive_mobs:
            c_x = mob.x+8
            c_y = mob.y+8
            if self.pos_covered_with_fog((c_x, c_y)):
                mutate_mobs.append(mob)

        random.shuffle(selected_nodes)

        for mob,node in zip(mutate_mobs, itertools.cycle(selected_nodes)):
            mob.do_fog_mutate(node.pos_px)

        self.for_maze.regenerate_selected(selected_nodes)




    def pos_covered_with_fog(self, pos):
        fog_x = pos[0]-self.x
        fog_y = pos[1]-self.y

        if fog_x < 0 or fog_x >= self.gamestate.maze.width_px:
            return False
        if fog_y < 0 or fog_y >= self.gamestate.maze.height_px:
            return False

        c = tuple(self.fog_surface.get_at((fog_x, fog_y)))
        return c[3] > 0

    def update(self):
        self.x += self.vx
        if self.x < -self.for_maze.width_px:
            self.passed = True

        if not self.mutated:
            self.mutate_count -= 1

            if self.mutate_count < 0:
                self.mutate()
                self.mutated = True
                self.surface = self.inverse_surface

            # don't flash too early.
            elif self.mutate_count < 200:
                self.flash_count += 10
                if self.flash_count > self.mutate_count:
                    self.flash_count = 0
                    if self.surface is self.inverse_surface:
                        self.surface = self.fog_surface
                    else:
                        self.surface = self.inverse_surface





    def render(self, to_surface):
        to_surface.blit(self.surface, (self.x,self.y))

