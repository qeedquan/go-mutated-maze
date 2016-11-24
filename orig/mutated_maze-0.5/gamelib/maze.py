__author__ = 'sonneveld'

import random
import pygame
import data
import ega
from math import hypot

MAZE_COLOUR = ega.BRIGHT_GREEN

class MazeNode(object):

    def __init__(self, id, x, y, maze_parent):
        self.id = id
        self.x = x
        self.y = y
        self.maze_parent = maze_parent

        self.colour = MAZE_COLOUR

        #start with a bunch of closed nodes
        self.up_open = False
        self.down_open = False
        self.left_open = False
        self.right_open = False

        self.is_rubble = False
        self.is_lava = False

    def set_wall(self, dir):
        if dir is 'up' and self.up_open:
            self.up_open = False
            other = self.up_node
            other.down_open = False
        elif dir is 'down' and self.down_open:
            self.down_open = False
            other = self.down_node
            other.up_open = False
        elif dir is 'left' and self.left_open:
            self.left_open = False
            other = self.left_node
            other.right_open = False
        elif dir is 'right' and self.right_open:
            self.right_open = False
            other = self.right_node
            other.left_open = False

    def clear_wall(self, dir):
        if dir is 'up' and not self.up_open and self.has_up_node:
            self.up_open = True
            other = self.up_node
            other.down_open = True
        elif dir is 'down' and not self.down_open and self.has_down_node:
            self.down_open = True
            other = self.down_node
            other.up_open = True
        elif dir is 'left' and not self.left_open and self.has_left_node:
            self.left_open = True
            other = self.left_node
            other.right_open = True
        elif dir is 'right' and not self.right_open and self.has_right_node:
            self.right_open = True
            other = self.right_node
            other.left_open = True

    def get_wall_state(self, dir):
        if dir is 'up': return self.up_open
        if dir is 'down': return self.down_open
        if dir is 'left': return self.left_open
        if dir is 'right': return self.right_open
        raise Exception('bad direction: '+dir)

    def get_node_in_dir(self, dir):
        if dir is 'up': return self.up_node
        if dir is 'down': return self.down_node
        if dir is 'left': return self.left_node
        if dir is 'right': return self.right_node
        raise Exception('bad direction: '+dir)

    def get_wall_count(self):
        count = 0
        if not self.up_open: count +=1
        if not self.down_open: count +=1
        if not self.left_open: count +=1
        if not self.right_open: count +=1
        return count

    @property
    def x_px(self):
        return self.x*16

    @property
    def y_px(self):
        return self.y*16

    @property
    def pos_px(self):
        return self.x*16,self.y*16

    def distance_from(self, node):
        return hypot(self.x-node.x, self.y-node.y)

    @property
    def hitbox(self):
        return pygame.Rect(self.x*16, self.y*16, 16, 16)

    @property
    def available_for_maze(self):
        return not self.is_rubble and not self.is_lava


    @property
    def open_nodes(self):
        """ nodes that have doors open to each other """
        if self.up_open: yield self.up_node
        if self.down_open: yield self.down_node
        if self.left_open: yield self.left_node
        if self.right_open: yield self.right_node

    @property
    def avail_directions(self):
        """ nodes that have doors open to each other """
        if self.up_open: yield 'up'
        if self.down_open: yield 'down'
        if self.left_open: yield 'left'
        if self.right_open: yield 'right'

    @property
    def nearby_nodes(self):
        """ nearby nodes that are cleared of rubble """
        if self.y > 0: yield self.up_node
        if self.y < self.maze_parent.height-1: yield self.down_node
        if self.x > 0: yield self.left_node
        if self.x < self.maze_parent.width-1: yield self.right_node

    @property
    def nearby_dirs(self):
        if self.y > 0: yield 'up'
        if self.y < self.maze_parent.height-1: yield 'down'
        if self.x > 0: yield 'left'
        if self.x < self.maze_parent.width-1: yield 'right'

    def close_all(self):
        if self.up_open:
            self.up_open = False
            self.up_node.down_open = False
        if self.down_open:
            self.down_open = False
            self.down_node.up_open = False
        if self.left_open:
            self.left_open = False
            self.left_node.right_open = False
        if self.right_open:
            self.right_open = False
            self.right_node.left_open = False

    def open_all(self):
        for dir in self.nearby_dirs:
            self.clear_wall(dir)

    @property
    def up_node(self):
        return self.maze_parent.get_node(self.x, self.y-1)
    @property
    def down_node(self):
        return self.maze_parent.get_node(self.x, self.y+1)
    @property
    def left_node(self):
        return self.maze_parent.get_node(self.x-1, self.y)
    @property
    def right_node(self):
        return self.maze_parent.get_node(self.x+1, self.y)

    @property
    def has_up_node(self):
        return self.y > 0
    @property
    def has_down_node(self):
        return self.y < self.maze_parent.height-1
    @property
    def has_left_node(self):
        return self.x > 0
    @property
    def has_right_node(self):
        return self.x < self.maze_parent.width-1


class Maze(object):

    def __init__(self, width, height):
        self.width = width
        self.height = height
        self._data = None

        self._reset_maze()

    @property
    def width_px(self):
        return self.width*16

    @property
    def height_px(self):
        return self.height*16

    def get_node(self, x, y):
        """ get the node at x, y
        :param x: x pos
        :param y: y pos
        :rtype: MazeNode
        """
        return self._data[y*self.width + x]

    def collides_nodes(self, rect):
        """ yield nodes that collide with rect """
        rect = pygame.Rect(rect)
        for node in self._data:
            if rect.colliderect(node.hitbox):
                yield node

    def render_to_surface(self, surface):
        """ render maze walls to surface """

        surface.fill((40,40,40,0xff))

        maze_width_px = self.width * 16
        maze_height_px = self.height*16

        pygame.draw.rect(surface, MAZE_COLOUR, (0,0,maze_width_px,maze_height_px), 1)

        for node in self._data:
            x,y = node.pos_px

            if not node.up_open:
                pygame.draw.line(surface, node.colour, (x, y), (x+15, y), 1)
            if not node.down_open:
                pygame.draw.line(surface, node.colour, (x, y+15), (x+15, y+15), 1)
            if not node.left_open:
                pygame.draw.line(surface, node.colour, (x, y), (x, y+15), 1)
            if not node.right_open:
                pygame.draw.line(surface, node.colour, (x+15, y), (x+15, y+15), 1)



    # GENERATION (KRUSKAL'S ALGORITHM)
    # ===============================================================

    def _kruskal_reset_maze(self):
        """ reset maze with new closed nodes     """
        self._data = []
        id = 0
        for y in xrange(self.height):
            for x in xrange(self.width):
                node = MazeNode(id, x, y, self)
                id += 1
                self._data.append(node)

    def kruskal_generate(self, nodes=None):
        """ simple maze generation        """

        if nodes is None:
            nodes = self._data

        while self._is_incomplete(nodes):

            # pick a random node
            x = random.randint(0, self.width-1)
            y = random.randint(0, self.height-1)
            node = self.get_node(x, y)
            if node not in nodes:
                continue
            if not node.available_for_maze:
                continue

            # pick random direction
            direction = random.randint(0, 3)
            if direction == 0: #UP
                if y == 0:
                    continue
                other = self.get_node(x, y-1)
                if other not in nodes:
                    continue
                if not other.available_for_maze:
                    continue
                if node.id == other.id:
                    continue
                node.up_open = True
                other.down_open = True
                self._set_nodes_to_id(other.id, node.id, nodes)

            elif direction == 1: #DOWN
                if y == self.height-1:
                    continue
                other = self.get_node(x, y+1)
                if other not in nodes:
                    continue
                if not other.available_for_maze:
                    continue
                if node.id == other.id:
                    continue
                node.down_open = True
                other.up_open = True
                self._set_nodes_to_id(other.id, node.id, nodes)


            elif direction == 2: #LEFT
                if x == 0:
                    continue
                other = self.get_node(x-1, y)
                if other not in nodes:
                    continue
                if not other.available_for_maze:
                    continue
                if node.id == other.id:
                    continue
                node.left_open = True
                other.right_open = True
                self._set_nodes_to_id(other.id, node.id, nodes)

            elif direction == 3: # RIGHT
                if x == self.width-1:
                    continue
                other = self.get_node(x+1, y)
                if other not in nodes:
                    continue
                if not other.available_for_maze:
                    continue
                if node.id == other.id:
                    continue
                node.right_open = True
                other.left_open = True
                self._set_nodes_to_id(other.id, node.id, nodes)


    def _is_incomplete(self, data):
        """ is the maze generation incomplete?
        :param data: list of maze nodes
        :return: true if not all ids are the same
        """
        prev = None
        for x in data:
            if x is None:
                return True
            if prev and x.id != prev.id:
                return True
            prev = x
        return False

    def _set_nodes_to_id(self, old_id, new_id, nodes=None):
        """
        set nodes from old_id to new_id
        :param old_id: old id value
        :param new_id: new id value
        """
        if nodes is None:
            nodes = self._data
        for x in nodes:
            if x.id == old_id:
                x.id = new_id

    # REGENERATION (KRUSKAL'S ALGORITHM)
    # ===============================================================

    def kruskal_regenerate_selected(self, selected):
        """ regenerate maze for selected nodes """

        # reset ids and close selected nodes
        #id = 0
        for node in self._data:
            node.id = None
            if node in selected:
                node.close_all()
                node.is_rubble = False

        # merge ids hopefully.
        id = 0
        for node in self._data:
            if node.id is not None:
                continue
            self._set_ids_from_start(node, id)
            id += 1


        node_sets = self._find_cleared_node_sets()
        for node_set in node_sets:
            self.kruskal_generate(node_set)


    def _set_ids_from_start(self, node, id):
        """ for node, if id is not None, set id then set the id for its
        open neighbours.   """
        if node.id is not None:
            return
        node.id = id

        if node.up_open:    self._set_ids_from_start(node.up_node, id)
        if node.down_open:  self._set_ids_from_start(node.down_node, id)
        if node.left_open:  self._set_ids_from_start(node.left_node, id)
        if node.right_open: self._set_ids_from_start(node.right_node, id)


    def _find_cleared_node_sets(self):
        """ find sets of nodes that are availble to walk through and near
        each other.  (eg, a trail of rubble could create two sets)
        :rtype: list of MazeNode
        """
        q = set(self._data)
        while q:
            n = q.pop()
            s = self._collect_nearby_cleared_nodes(n)
            if s:
                yield s
                q -= s


    def _collect_nearby_cleared_nodes(self, node, s=None):
        """
        :type s: set
        :type node: MazeNode
        :rtype: set
        """
        if s is None:
            s = set()
        if node.is_rubble:
            return

        if node in s:
            return
        s.add(node)
        for x in node.nearby_nodes:
            self._collect_nearby_cleared_nodes(x, s)
        return s



    # GENERATION (BRAID)
    # ===============================================================

    '''
    http://www.astrolog.org/labyrnth/algrithm.htm
    Braid: To create a Maze without dead ends, basically add wall segments
    throughout the Maze at random, but ensure that each new segment added will
    not cause a dead end to be made. I make them with four steps: (1) Start with
    the outer wall, (2) Loop through the Maze and add single wall segments
    touching each wall vertex to ensure there are no open rooms or small "pole"
    walls in the Maze, (3) Loop over all possible wall segments in random
    order, adding a wall there if it wouldn't cause a dead end, (4) Either run
    the isolation remover utility at the end to make a legal Maze that has a
    solution, or be smarter in step three and make sure a wall is only added
    if it also wouldn't cause an isolated section.
    '''

    def _braid_reset_maze(self):
        self._data = []
        id = 0
        for y in xrange(self.height):
            for x in xrange(self.width):
                node = MazeNode(id, x, y, self)
                node.up_open = y != 0
                node.down_open = y != self.height-1
                node.left_open = x != 0
                node.right_open = x != self.width-1
                id += 1
                self._data.append(node)

    def braid_generate(self, nodes=None):

        if nodes is None:
            nodes = self._data

        # generate a list of wall segments
        walls = []
        for node in nodes:
            for dir in node.avail_directions:
                w = (node, dir)
                walls.append(w)

        random.shuffle(walls)

        for node,dir in walls:
            if not node.get_wall_state(dir):
                continue
            if node.get_wall_count() >= 2:
                continue
            other = node.get_node_in_dir(dir)
            if other.get_wall_count() >= 2:
                continue
            if not self._is_connected(node,other):
                continue
            node.set_wall(dir)


    def _is_connected(self, first, second):
        ''' can you get from first to the second without going through the
          nearby connection '''

        seen = set([first, second])
        queue = set([x for x in first.open_nodes if x not in (first, second)])

        while queue:
            n = queue.pop()
            seen.add(n)
            for x in n.open_nodes:
                if x is second:
                    return True
                if x in seen:
                    continue
                queue.add(x)

        return False


    def braid_regenerate_selected(self, selected):
        for node in selected:
            node.open_all()
        self.braid_generate(selected)


    # DEFAULTS
    # ===============================================================

    # set these to set the default algorithm
    #_reset_maze = _kruskal_reset_maze
    #generate = kruskal_generate
    #regenerate_selected = kruskal_regenerate_selected
    _reset_maze = _braid_reset_maze
    generate = braid_generate
    regenerate_selected = braid_regenerate_selected
