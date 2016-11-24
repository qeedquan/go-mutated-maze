__author__ = 'sonneveld'

import pygame

import ega

class TextBufferRenderer(object):


    def __init__(self, fontname, size, default_colour=ega.BRIGHT_MAGENTA):
        self.font  = pygame.font.SysFont(fontname, size)
        self.default_colour = default_colour
        self.clear()

    def clear(self):
        self.text_list = []

    def add(self, text, colour=None):
        if colour is None:
            colour = self.default_colour
        self.text_list.append( (text, colour))


    def render_to_surface(self, surface, pos):
        x, y = pos
        for (text, colour) in self.text_list:
            s = self.font.render(text, True, colour)
            surface.blit(s, (x,y))
            y += s.get_height()
