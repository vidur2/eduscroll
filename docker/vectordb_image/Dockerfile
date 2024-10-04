FROM ubuntu
RUN apt-get update
RUN apt-get -y install python3-pip
WORKDIR /app
RUN pip install chromadb
EXPOSE 8000
CMD ["chroma", "run", "--host", "0.0.0.0"]