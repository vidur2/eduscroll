from torch import nn
from minute_model import MinuteModel

class ReccomendationLoss(nn.Module):
    def __init__(self, d_model=16, n_head=4):
        super(nn.Module, self).__init__()
        self.model = MinuteModel(d_model, n_head)
    
    def forward(self, state, action):
        cpy = [i for i in state] + [action]
        return self.model(cpy)[0] - self.model(state)[0]