FROM ubuntu
RUN apt-get update
RUN apt-get -y install python3-pip
RUN apt-get install -y ffmpeg
# set work directory
WORKDIR /app
# set environment variables
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

# install dependencies
RUN pip install --upgrade pip
COPY requirements.txt /app
RUN pip install -r requirements.txt
EXPOSE 8003

# copy project
COPY . /app

CMD ["gunicorn", "--bind", ":8003", "--workers", "3", "coreapi.wsgi:application", "--timeout", "6000"]
