import sys
import os
import random
import pygame
import pygame.image

import ega
import data


def set_colour(surface, colour):
    p = pygame.PixelArray(surface)
    p.replace(ega.BLACK, colour, 0.0)

def set_random_colour(surface, colour_range):
    colour = random.choice(colour_range)
    set_colour(surface, colour)


def load_sprite_surface(path, colour_range=None, force_colour=None):
    p = data.filepath(path)
    i = pygame.image.load(p)

    if force_colour:
        set_colour(i, force_colour)
    elif colour_range:
        set_random_colour(i, colour_range)
        
    i = i.convert_alpha()
    return i