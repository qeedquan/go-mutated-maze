import sys
import os
import random
import pygame
import pygame.image


# unused giant sprite


class Giant(object):


    def __init__(self, gamestate):

        self.gamestate = gamestate
        self.for_maze = gamestate.maze

        self.scale = random.randint(3, 10)
        self.surface = self.generate_surface(self.scale)
        self.w, self.h = self.surface.get_size()


        startpos = random.randint(0, 3)
        startpos = 0
        if startpos == 0:  # from the left
            self.x = 0 - self.w
            self.y = random.randint(0, self.for_maze.height_px - self.h)
            self.vx = 3
            self.vy = 0




            #self.vx = 0
            #self.vy = 0
            #self.scale = random.randint(3, 10)
            #self.surface = self.generate_surface(self.scale)



    def generate_surface(self, scale):
        giant_path = data.filepath("bad_person.png")
        giant_image = pygame.image.load(giant_path)
        set_colour(giant_image, None)
        size = giant_image.get_size()
        size = (size[0]*scale, size[1]*scale)

        giant_image = pygame.transform.scale(giant_image, size)
        giant_image = giant_image.convert_alpha()
        return giant_image

    def update(self):
        self.x += self.vx
        self.y += self.vy

        if self.x >= self.for_maze.width_px:
            return
        if self.y >= self.for_maze.height_px:
            return

        for node in self.for_maze.collides_nodes((self.x, self.y, self.w, self.h)):
            if not node.is_rubble:
                node.is_rubble = True
                node.close_all()



    def render(self, to_surface):
        to_surface.blit(self.surface, (self.x, self.y))



