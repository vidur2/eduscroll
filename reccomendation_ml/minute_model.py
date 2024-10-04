import torch.optim as optim
import torch
import math
from torch import nn, Tensor
from torch.nn import TransformerEncoder, TransformerEncoderLayer, Module, Linear
import random

class PositionalEncoding(nn.Module):

    def __init__(self, d_model: int, dropout: float = 0.1, max_len: int = 5000):
        super().__init__()
        self.dropout = nn.Dropout(p=dropout)

        position = torch.arange(max_len).unsqueeze(1)
        div_term = torch.exp(torch.arange(0, d_model, 2) * (-math.log(10000.0) / d_model))
        pe = torch.zeros(max_len, 1, d_model)
        pe[:, 0, 0::2] = torch.sin(position * div_term)
        pe[:, 0, 1::2] = torch.cos(position * div_term)
        self.register_buffer('pe', pe)

    def forward(self, x: Tensor) -> Tensor:
        """
        Arguments:
            x: Tensor, shape ``[seq_len, batch_size, embedding_dim]``
        """
        x = x + self.pe[:x.size(0)]
        return self.dropout(x)

class MinuteModel(Module):
    def __init__(self, d_model=16, nhead=4, dropout=0.5):
        super().__init__()
        self.pos_encoder = PositionalEncoding(d_model, dropout)
        encoder_layer = TransformerEncoderLayer(d_model, nhead)
        self.transformer_encoder = TransformerEncoder(encoder_layer, num_layers=4)
        self.linear = Linear(d_model, 1)
        self.init_weights()
        self.d_model = d_model
    
    def init_weights(self) -> None:
        initrange = 0.1
        self.linear.bias.data.zero_()
        self.linear.weight.data.uniform_(-initrange, initrange)

    def forward(self, src: Tensor, src_mask: Tensor = None, device="cpu") -> Tensor:
        """
        Arguments:
            src: Tensor, shape ``[seq_len, batch_size]``
            src_mask: Tensor, shape ``[seq_len, seq_len]``

        Returns:
            output Tensor of shape ``[seq_len, batch_size, ntoken]``
        """
        src = self.pos_encoder(src)
        if src_mask is None:
            """Generate a square causal mask for the sequence. The masked positions are filled with float('-inf').
            Unmasked positions are filled with float(0.0).
            """
            src_mask = nn.Transformer.generate_square_subsequent_mask(len(src)).to(device)
        output = self.transformer_encoder(src, src_mask)
        output = self.linear(output)
        return output[0][0]
    
    def train(self, data, iterations):
        criterion = nn.MSELoss()
        optimizer = optim.SGD(self.parameters(), lr=0.001, momentum=0.9)
        running_loss = 0.0
        for epoch in range(iterations):
            for i, point in enumerate(data, 0):
                inputs, labels = point
                optimizer.zero_grad()

                outputs = self(inputs)
                print(outputs)
                loss = criterion(outputs, labels)
                loss.backward()
                optimizer.step()

                running_loss += loss.item()
                if i % 1000 == 999:    # print every 1000 mini-batches
                    print(f'[{epoch + 1}, {i + 1:5d}] loss: {running_loss / 2000:.3f}')
                    running_loss = 0.0
        print("Finished training")

if (__name__ == "__main__"):
    model = MinuteModel(d_model=1, nhead=1)
    model.train([(torch.tensor([[random.random() * 100] for _ in range(random.randint(1, 10))]), torch.tensor([random.random() * 100])) for _ in range(1000)], 100)