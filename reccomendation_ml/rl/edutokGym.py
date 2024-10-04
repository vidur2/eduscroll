from gym import Env
import sys
sys.path.append("..")
from minute_model import MinuteModel
import random
import numpy as np
import torch
from torch import Tensor

class EdutokEnvironment(Env):
    def __init__(self):
        self.model = MinuteModel(4, 1)
        super(Env, self).__init__()
        self.reward = 0 # minutes spent
        self.state = [torch.tensor([random.random() * 100 - 50 for i in range(self.model.d_model)])] # previous state
    
    def reset(self, seed = None):
        random.seed(seed)
        self.reward = 0
        self.state = [[random.random() * 100 - 50 for i in range(self.model.d_model)]]
        return np.array(self.state), self.reward
    
    def step(self, action):
        self.state.append(action)
        reward = self.model(torch.tensor(self.state))
        out = []
        for i in self.state:
            if (type(i) == Tensor):
                out.append(i.detach().numpy())
            else:
                out.append(np.array(i))
        return np.array(out), reward, reward < 0, None, None