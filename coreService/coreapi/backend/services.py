import pathlib
import chromadb
from chromadb.utils import embedding_functions
from sentence_transformers import SentenceTransformer
from botocore.exceptions import ClientError
import os
from langchain.prompts import PromptTemplate
from langchain.llms import Bedrock
from langchain.schema import HumanMessage, SystemMessage
from langchain.chat_models import BedrockChat
from langchain.prompts.chat import (
    ChatPromptTemplate,
    SystemMessagePromptTemplate,
    HumanMessagePromptTemplate,
    MessagesPlaceholder
)
from langchain.chains import LLMChain

import requests
from json import dumps, loads
import boto3
from PIL import Image
from io import BytesIO
from django.http import HttpRequest
import base64
from uuid import uuid4
import ffmpeg
from random import randint
from PIL import Image
import shutil
import time
from langchain.memory import ConversationBufferMemory
import subprocess
import dotenv


modelPath = "model"

CHUNK_SIZE = 1024

if (__name__  == "__main__"):
    modelPath = "../model"

model = SentenceTransformer(modelPath)
s3Resource = boto3.resource('s3', aws_access_key_id="API_KEY_HERE",aws_secret_access_key="API_KEY_HERE")
s3Client = boto3.client('s3', region_name = "us-east-1", aws_access_key_id="API_KEY_HERE",aws_secret_access_key="API_KEY_HERE")
bedrock_runtime = boto3.client(service_name="bedrock-runtime", region_name = "us-east-1", aws_access_key_id="API_KEY_HERE",aws_secret_access_key="API_KEY_HERE")
polly = boto3.client('polly', region_name = "us-east-1", aws_access_key_id="API_KEY_HERE",aws_secret_access_key="API_KEY_HERE")
voices = ['Aditi','Amy','Astrid','Bianca','Brian','Camila','Carla','Carmen','Celine','Chantal','Conchita','Cristiano','Dora','Emma','Enrique','Ewa','Filiz','Gabrielle','Geraint','Giorgio','Gwyneth','Hans','Ines','Ivy','Jacek','Jan','Joanna','Joey','Justin','Karl','Kendra','Kevin','Kimberly','Lea','Liv','Lotte','Lucia','Lupe','Mads','Maja','Marlene','Mathieu','Matthew','Maxim','Mia','Miguel','Mizuki','Naja','Nicole','Olivia','Penelope','Raveena','Ricardo','Ruben','Russell','Salli','Seoyeon','Takumi','Tatyana','Vicki','Vitoria','Zeina','Zhiyu','Aria','Ayanda','Arlet','Hannah','Joanna','Daniel','Liam','Pedro','Kajal','Hiujin','Laura','Elin','Ida','Suvi','Ola','Hala','Andres','Sergio','Remi','Adriano','Thiago','Ruth','Stephen','Kazuha','Tomoko','Niamh','Sofie','Lisa','Isabelle','Zayd','Danielle','Gregory']
intros = [
    'Have you ever wondered', 
    'Interested in finding out', 
    'Curious to explore', 
    'Ready to delve into', 
    'Prepare to uncover', 
    'Ever pondered about', 
    'Did you ever think about', 
    'Eager to discover', 
    'Excited to learn about', 
    'Wondering what lies behind', 
    'Ever considered', 
    'Keen to investigate', 
    'Have you ever pondered', 
    'Intrigued to know more', 
    'Fascinated by the idea of', 
    'Dying to find out', 
    'Want to unravel', 
    'Have you ever reflected on', 
    'Interested to explore', 
    'Are you ready to discover', 
    'Have you ever contemplated', 
    'Hoping to uncover', 
    'Ever asked yourself', 
    'Eager to delve into', 
    'Ever thought about', 
    'Enthusiastic to learn about', 
    'Intrigued by the concept of', 
    'Ever considered the possibility of', 
    'Curious to know more about', 
    'Ready to explore', 
    'Have you ever mused over', 
    'Did you ever wonder what', 
    'Fascinated to explore', 
    'Excited to delve into', 
    'Ever imagined', 
    'Pondering the mysteries of', 
    'Are you curious about', 
    'Wondering about the secrets of', 
    'Intrigued to discover', 
    'Are you ready to uncover', 
    'Ever contemplated', 
    'Eager to know more about', 
    'Ready to investigate', 
    'Have you ever questioned', 
    'Interested in delving into', 
    'Are you curious to learn more about', 
    'Are you intrigued by', 
    'Have you ever explored the idea of', 
    'Ever mused over', 
    'Wondering what might be', 
    'Did you ever consider', 
    'Interested in the mysteries of', 
    'Have you ever imagined what', 
    'Eager to uncover the truth about', 
    'Ever pondered the idea of', 
    'Excited to dive into', 
    'Curious to know what', 
    'Have you ever reflected upon', 
    'Are you keen to explore', 
    'Ready to ponder', 
    'Ever questioned why', 
    'Intrigued by the possibilities of', 
    'Curious about the origins of', 
    'Did you ever think to yourself', 
    'Interested in uncovering', 
    'Fascinated by the mysteries of', 
    'Are you eager to learn about', 
    'Have you ever contemplated the idea of', 
    'Ever pondered over the thought of', 
    'Curious to discover', 
    'Did you ever ponder what', 
    'Eager to explore the world of', 
    'Interested in discovering', 
    'Ever thought about exploring', 
    'Fascinated by the concept of', 
    'Are you ready to ponder', 
    'Wondering what the future holds for', 
    'Intrigued to learn more about', 
    'Have you ever considered exploring', 
    'Curious to find out more about', 
    'Are you fascinated by', 
    'Ready to contemplate', 
    'Ever wondered why', 
    'Interested in delving deeper into', 
    'Have you ever been curious about', 
    'Did you ever ask yourself', 
    'Eager to delve deeper into', 
    'Excited to explore the world of', 
    'Ready to ponder upon', 
    'Ever pondered the mysteries of', 
    'Curious to know the truth about', 
    'Have you ever pondered over', 
    'Fascinated to learn more about', 
    'Wondering about the significance of', 
    'Interested in exploring the depths of', 
    'Ever considered the implications of', 
    'Are you ready to contemplate the idea of', 
    'Have you ever wondered what it would be like to'
]


class SignedUrlStore:
    def __init__(self):
        self.store = {}
    def get_signed_url(self, url, bucket_name= "eduscroll-video-output"):
        if (url in self.store and not self.store[url].isExpired()):
            return self.store[url].getUrl()
        splitSlash = url.split("/")
        obj = splitSlash[len(splitSlash) - 1]
        obj = obj.replace("\\", "").replace('"',"")
        out = s3Client.generate_presigned_url(
            ClientMethod='get_object',
            Params={
                'Bucket': bucket_name,
                'Key': obj
            },
            ExpiresIn=3600 # one hour in seconds, increase if needed
        )
        self.store[url] = SignedUrl(out)
        return out

class SignedUrl:
    def __init__(self, url, expireTime = 3600):
        self.url = url
        self.expire = time.time() + expireTime - 300

    def getUrl(self):
        if (not self.isExpired()):
            return self.url
        else:
            raise Exception("Url is expired")

    def isExpired(self):
        return time.time() > self.expire

def chat_generate_problem(subject, ctx):
    llm = BedrockChat(model_id="meta.llama2-70b-chat-v1", client = bedrock_runtime)
    hostname = "edutok-chroma-1"
    # client = chromadb.HttpClient(host=hostname, port = 8000)

    sys_prompt =  """
        You are a machine that generates 100-word educational videos on {school_subject} based on a conversation with an educator. 
        Use the following previously written questions as references for creating new material.
        Here is a passage from which you can generate videos: {context1}
        
        Now, create a new problem on your subject based on the problems that you have already referenced. Don't tell me what the problem is.  I'm going to give you some tasks to do
        based on the problem you've created.

        Previous conversation: 
        {chat_history}

        Educator Instructions: {request}
        """
        # input_variables=["school_subject", "context1", "request"]
   
    memory = ConversationBufferMemory(memory_key="chat_history", return_messages=True)
    chat_prompt = ChatPromptTemplate(
        messages = 
        [
            SystemMessagePromptTemplate.from_template(template = sys_prompt, partial_variables={"context1" : ctx, "school_subject" : subject}),
            MessagesPlaceholder(variable_name="chat_history"), 
            HumanMessagePromptTemplate.from_template("{request}")
        ]
    )


    # collection = client.get_collection(subject)
    # collection_query = collection.query(
    #     query_texts=["Give me some questions about " + subject],
    #     n_results = 1,
    # )
    # print(collection_query["documents"][0][0])
    conversation = LLMChain(
        llm = llm,
        prompt = chat_prompt,
        verbose = True,
        memory = memory

    )
    # req = HttpRequest()
    # res = requests.post(req.build_absolute_uri("/query_cache/add"), data=dumps({
    #     "docs": [script], 
    #     "subject": subject,
    #     "problems": [{"question": 'TEST', "answer": ""}],
    #     "s3VideoUri": ["s3://test"]
    # }))
    return None

def chat_generate_script(subject,ctx):
    llm = BedrockChat(model_id="meta.llama2-70b-chat-v1", client = bedrock_runtime)
    hostname = "edutok-chroma-1"
    if (__name__ == "__main__"):
        hostname = "localhost"
    # client = chromadb.HttpClient(host=hostname, port = 8000)

    sys_prompt =  """
        You are a machine that generates 30-second educational videos on {school_subject} based on a conversation with an educator. 
        Use the following previously written questions as references for creating new material.
        Here is the textbook page which you are to use as context: {context1}
        
        Now, create a new problem on your subject based on the problems that you have already referenced. Don't tell me what the problem is.  I'm going to give you some tasks to do
        based on the problem you've created.

        Previous conversation: 
        {chat_history}

        Educator Instructions: {request}
        """
        # input_variables=["school_subject", "context1", "request"]

    script_prompt_1 =  f"""
    Now we are going to use the problem that you created to make a complete 30-second video script. For right now, you need to make the first 15 seconds of the videoby introducing and describing the problem in a way that is accessible for high school students.
    The script should use clear language and explain the premise, background, and solution of the problem effectively. Don't include stage directions or suggestions for visuals or audio. Also don't tell me when its complete or give me any other information outside of just the script. Also,
    don't say "Sure! Here's the script for the first 15 seconds of the video:" or anything similar. Make sure it is short form. Introduce the script with: {intros[randint(0, len(intros) - 1)]}. Limit your response to 10 words.
    """
    script_prompt_2 = """
    Now that the first half of the script is complete, write the last 15 seconds of the script. Based on the first 15 seconds of the script, make sure to thoroughly explain and walk through the problem so that our video can be complete. Write it in the same style as before.
    Make sure the script has an sensible conclusion and then don't write anything else. Also don't tell me when its complete or give me any other information outside of just the script.Don't include stage directions or suggestions for visuals or audio. Also don't tell me when its complete or give me any other information outside of just the script. Also,
    don't say "Sure! Here's the script for the first 15 seconds of the video:" or anything similar. Limit your response to 10 words.
    """
   
    memory = ConversationBufferMemory(memory_key="chat_history", return_messages=True)
    chat_prompt = ChatPromptTemplate(
        messages = 
        [
            SystemMessagePromptTemplate.from_template(template = sys_prompt, partial_variables={"context1" : ctx, "school_subject" : subject}),
            MessagesPlaceholder(variable_name="chat_history"), 
            HumanMessagePromptTemplate.from_template("{request}")
        ]
    )


    # collection = client.get_collection(subject)
    # collection_query = collection.query(
    #     query_texts=["Give me some questions about " + subject],
    #     n_results = 1,
    # )
    # print(collection_query["documents"][0][0])
    conversation = LLMChain(
        llm = llm,
        prompt = chat_prompt,
        verbose = False,
        memory = memory

    )

    # script = conversation({"request": script_prompt_1}) + "\n" + conversation({"request": script_prompt_2})
    # req = HttpRequest()
    # res = requests.post(req.build_absolute_uri("/query_cache/add"), data=dumps({
    #     "docs": [script], 
    #     "subject": subject,
    #     "problems": [{"question": 'TEST', "answer": ""}],
    #     "s3VideoUri": ["s3://test"]
    # }))
    p1 = conversation({"request": script_prompt_1})['text'] 
    p2 = conversation({"request": script_prompt_2})['text']
    
    return shorten_prompt_iter(trim(p1) + "\n\n" + trim(p2), 115, 5, llm)

def trim(script):
    if ("15" in script.lower() or "30" in script.lower()):
        splitScript = script.split("\n\n")
        return "\n".join(splitScript[1:len(splitScript) - 1])
    else:
        return script

def generate_problem(subject):
 
    llm = BedrockChat(model_id="meta.llama2-70b-chat-v1", client = bedrock_runtime)
    hostname = "edutok-chroma-1"
    if (__name__ == "__main__"):
        hostname = "localhost"
    client = chromadb.HttpClient(host=hostname, port = 8000)

    sys_prompt = PromptTemplate(
        template = 
        """
        You are a machine that generates 100-word educational videos on {school_subject}. 
        Use the following previously written questions as references for creating new material.
        This is the first example: {context_1}
        """,
        input_variables=["school_subject", "context_1"]
    )
    problem_prompt = PromptTemplate(
        template = 
        """
        First, create a new problem on your subject based on the problems that you have already referenced. Don't tell me what the problem is. Now I'm going to give you some tasks to do
        based on the problem you've created.
        """,
        input_variables=[]
    )
    script_prompt = PromptTemplate(
        template = 
        """
        First use the problem that you just created to write a complete spoken script for an approximately 50-word educational video describing the problem in a way that is accessible for high school students.
        The script should use clear textand explain the premise and background of the problem effectively. Make sure the script has an sensible conclusion and then don't write anything else. Don't include stage directions or suggestions for visuals or audio. 
        Also don't tell me when its complete or give me any other information outside of just the script.
        """,
        input_variables=[]
    )
    system_prompt = SystemMessagePromptTemplate(prompt = sys_prompt)
    user_prompt = HumanMessagePromptTemplate(prompt = problem_prompt)
    content_prompt = HumanMessagePromptTemplate(prompt = script_prompt)
    chat_prompt = ChatPromptTemplate.from_messages([system_prompt, user_prompt, content_prompt])


    collection = client.get_collection(subject)
    collection_query = collection.query(
        query_texts=["Give me some questions about " + subject],
        n_results = 1,
    )
    print(collection_query["documents"][0][0])
    prompt_values = chat_prompt.format_prompt( 
                                                school_subject = subject, \
                                                context1=collection_query["documents"][0][0], \
    )
    
    script = llm(prompt_values.to_messages()).content
    # req = HttpRequest()222
    # res = requests.post(req.build_absolute_uri("/query_cache/add"), data=dumps({
    #     "docs": [script], 
    #     "subject": subject,
    #     "problems": [{"question": 'TEST', "answer": ""}],
    #     "s3VideoUri": ["s3://test"]
    # }))
    return script

def shorten_prompt_iter(orig: str, max_words: int, max_iter: int, llm):
    curr_words = len(orig.split(" "))
    i = 0
    sys_prompt = """
    You are a machine which makes scripts for tik toks shorter, given a script. Make sure to make it is still a cohesive script and finishes. Only include the script, no extra words, and no emojis.
    The script which you are shortening is: {script}
    """
    chat_prompt = ChatPromptTemplate(
        messages=
        [
            SystemMessagePromptTemplate.from_template(template=sys_prompt),
            HumanMessagePromptTemplate.from_template("{script}")
        ]
    )

    conversation = LLMChain(
        llm = llm,
        prompt = chat_prompt,
        verbose = False,
    )

    while (curr_words > max_words and i < max_iter):
        orig = conversation({"script": orig})['text'] 
        if ('"' in orig):
            orig = orig.split('"')[1]
        tmp = orig.split("\n\n")
        if ("shorter" in tmp[0] or "shorter" in tmp[0]):
            tmp = tmp[1:]
            orig = "\n".join(tmp)
        curr_words = len(orig.split(" "))
        i += 1
    return orig

def gen_video(mapping, uuid):
    ffmpegInputs = []
    ffmpegAudio = []
    for i, av in enumerate(mapping):
        image, audio = av
        os.makedirs("tmpVid", exist_ok=True)
        os.makedirs("finalVid", exist_ok=True)
        os.system(f"ffmpeg -loop 1 -i {image} -i {audio} -c:v libx264 -tune stillimage -c:a aac -b:a 192k -pix_fmt yuv420p -shortest tmpVid/{uuid}-{i}.mp4")
        ffmpegInputs.append(ffmpeg.input(f"tmpVid/{uuid}-{i}.mp4"))
        ffmpegAudio.append(ffmpeg.input(audio))
    ffmpeg.concat(*ffmpegAudio, v =0, a=1).output(f"finalVid/{uuid}-audio.mp3").run(overwrite_output=True)
    ffmpeg.concat(*ffmpegInputs).output(f"finalVid/{uuid}-final.mp4").run(overwrite_output=True)
    # ffmpeg.concat(ffmpeg.input(f"finalVid/{uuid}-audio.mp3"), ffmpeg.input(f"finalVid/{uuid}-final.mp4")).output("finalVid/final.mp4").run(overwrite_output=True)
    # (ffmpeg.output(ffmpeg.input(f"finalVid/{uuid}-audio.mp3"), ffmpeg.input(f"finalVid/{uuid}-final.mp4"),f"finalVid/{uuid}-final-audio.mp4", codec='copy', vcodec='copy', acodec='aac').run(overwrite_output=True))
    os.system(f"ffmpeg -i finalVid/{uuid}-final.mp4 -i finalVid/{uuid}-audio.mp3 -map 0:v -map 1:a -c:v copy -shortest finalVid/{uuid}-combined.mp4")
    os.system(f"ffmpeg -i finalVid/{uuid}-combined.mp4 -vf ass=./audio/{uuid}.ass finalVid/{uuid}-combined-ass.mp4")
    shutil.rmtree("tmpVid")
    shutil.rmtree("audio")
    shutil.rmtree("images")
    return f"finalVid/{uuid}-combined-ass.mp4"



def gen_content(prompt_response, uuid):
    splitPrompt = prompt_response.split(".")
    out = []
    j = 0
    srtInfo = []
    paths = []
    for i in range(0, len(splitPrompt), len(splitPrompt) // 4):
        curr = " ".join(splitPrompt[i:i + (len(splitPrompt)) // 4]) if i < 7 * len(splitPrompt) // 4 else " ".join(splitPrompt[i:])
        audioPath, outFile = gen_audio(curr, uuid, j)
        paths.append(audioPath)
        srtInfo.append(outFile)
        out.append((gen_images(curr, uuid, i), audioPath))
        j += 1
    gen_srt(srtInfo, paths, uuid)
    return out

def generate_image_prompt(subprompt):
    llm = BedrockChat(model_id="meta.llama2-70b-chat-v1", client = bedrock_runtime)
    hostname = "edutok-chroma-1"
    if (__name__ == "__main__"):
        hostname = "localhost"
    client = chromadb.HttpClient(host=hostname, port = 8000)

    sys_prompt = PromptTemplate(
        template = 
        """
        You are a machine which generates image prompts from fragments of textbook passages. These prompts will be fed into stable diffusion. Just return the prompt, nothing else. Use words which are not specific to the context. Expand on what you are describing.
        """,
        input_variables=[]
    )
    script_prompt = PromptTemplate(
        template = 
        """
        The fragment which you are generating a prompt for is: {prompt}
        """,
        input_variables=["prompt"]
    )
    system_prompt = SystemMessagePromptTemplate(prompt = sys_prompt)
    content_prompt = HumanMessagePromptTemplate(prompt = script_prompt)
    chat_prompt = ChatPromptTemplate.from_messages([system_prompt, content_prompt])
    prompt_values = chat_prompt.format_prompt(prompt=subprompt)
    llm = BedrockChat(model_id="meta.llama2-70b-chat-v1", client = bedrock_runtime)
    return llm(prompt_values.to_messages()).content.replace("Prompt: ", "")

def gen_images(subprompt, uuid, i):
    subprompt = generate_image_prompt(subprompt)
    print(subprompt)
    tmp = bedrock_runtime.invoke_model(
                    body=dumps({
                        "text_prompts": [
                        {
                            "text": subprompt[:1200],
                            "weight": 0.7
                        }, 
                        {
                            "text": "DO NOT INCLUDE TEXT, OR LETTERS, OR ANYTHING WRITTEN",
                            "weight": 0.3
                        }, 
                        {
                            "text": "lowres, text, typography, letters", 
                            "weight": -1
                        },
                    ],
                        "samples" : 1,
                        "width": 640,
                        "height": 1536,                        
                    }),
                    modelId="stability.stable-diffusion-xl-v1",
                    accept="application/json", 
                    contentType="application/json"
                )
    os.makedirs("images", exist_ok=True)
    response_body = loads(tmp.get("body").read())
    base_64_img_str = response_body["artifacts"][0]["base64"]
    image = Image.open(BytesIO(base64.decodebytes(bytes(base_64_img_str, "utf-8"))))
    image.save(f"images/{uuid}-{i}.png")
    return f"images/{uuid}-{i}.png"

def gen_audio(subprompt, uuid, i):
    eleven_headers = {
        "Accept": "audio/mpeg",
        "Content-Type": "application/json",
        "xi-api-key": os.getenv("ELEVEN_API_KEY")
    }
    data = {
        "text": subprompt,
        "model_id": "eleven_monolingual_v1",
        "voice_settings": {
            "stability": 0.5,
            "similarity_boost": 0.5
        }
    }

    response = requests.post("https://api.elevenlabs.io/v1/text-to-speech/T5cu6IU92Krx4mh43osx", json=data, headers=eleven_headers)
    if (response.status_code != 200):
        print(response.text)
    os.makedirs("audio", exist_ok=True)
    with open(f"audio/{uuid}-{i}.mp3", 'wb') as f:
        for chunk in response.iter_content(chunk_size=CHUNK_SIZE):
            if chunk:
                f.write(chunk)
    res = polly.synthesize_speech(
        Engine='neural',
        LanguageCode='en-US',
        OutputFormat='json',
        Text=subprompt,
        VoiceId="Arthur",
        SpeechMarkTypes=["word"]
    )
    tmp = res['AudioStream'].read().decode("utf-8").split("\n")
    outFile = [loads(i) for i in tmp[:len(tmp) - 2]]
    return f"audio/{uuid}-{i}.mp3", outFile

def convertTotalToString(total):
    hours = convertHMToString(total // (1000 * 60 * 60))
    minutes = convertHMToString(total // (1000 * 60))
    ms = convertMsToString(total % (1000 * 60))
    ms = ms[:2] + "," + ms[2:]
    return hours, minutes, ms

def get_length(input_video):
    result = subprocess.run(['ffprobe', '-v', 'error', '-show_entries', 'format=duration', '-of', 'default=noprint_wrappers=1:nokey=1', input_video], stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    return round(float(result.stdout) * 1000)


def gen_srt(pollyOutput, audioPath, uuid, wordWindow=3):
    outStr = ""
    prevLastTime = 0
    nextTime = prevLastTime
    for j, speechmarks in enumerate(pollyOutput):
        for i in range(0, len(speechmarks), wordWindow):
            nxt = ""
            if (i < len(speechmarks) - wordWindow):
                hours, minutes, ms = convertTotalToString(prevLastTime + speechmarks[i + wordWindow]["time"])
                nxt = f"{hours}:{minutes}:{ms}"
            else:
                totalMs = get_length(audioPath[j]) + prevLastTime
                nextTime = totalMs
                hours, minutes, ms = convertTotalToString(totalMs)
                nxt = f"{hours}:{minutes}:{ms}"
            total = prevLastTime + int(speechmarks[i]["time"])
            hours, minutes, ms = convertTotalToString(total)
            if (j != 0):
                outStr += "\n"
            values = ""
            for k in speechmarks[i:min(i+wordWindow, len(speechmarks))]:
                values += f"{k['value']} "
            timeStr = f"{hours}:{minutes}:{ms} --> {nxt}"
            outStr += f"\n{i + 1}\n{timeStr}\n{values}"
        print(audioPath[j])
        prevLastTime = get_length(audioPath[j]) + prevLastTime
    with open(f"./audio/{uuid}.srt", "w") as f:
        f.write(outStr)
    os.system(f"ffmpeg -i ./audio/{uuid}.srt ./audio/{uuid}.ass")
    with open(f"./audio/{uuid}.ass", "r") as f:
        content = f.read().split("\n")
        style = content[10]
    style = style.split(",")
    style[len(style) - 2] = "155"
    style = ",".join(style)
    content[10] = style
    with open(f"./audio/{uuid}.ass", "w") as f:
        f.write("\n".join(content))


def convertHMToString(integerValue):
    if (integerValue < 10):
        return "0" + str(integerValue)
    else:
        return str(integerValue)

def convertMsToString(integerValue):
    val = str(integerValue)
    for i in range(len(val),5):
        val = "0" + val
    return val

def uploadToS3AndDelete(path):
    name = "eduscroll-video-output"
    try:
        s3Resource.meta.client.head_bucket(Bucket=name)
    except ClientError:
        s3Client.create_bucket(
            Bucket=name
        )
    
    with open(path, 'rb') as data:
        response = s3Client.upload_fileobj(data, name, path.split("/")[1])
    print(response)
    os.remove(path)
    return f"https://s3-us-east-1.amazonaws.com/{name}/{path.split('/')[1]}"

class DummyRequest:
    def __init__(self, data, request):
        self.data = data
        self.request = request
    def build_absolute_uri(self, path):
        return self.request.build_absolute_uri(path)

test_prompt = """In the past, biologists grouped living organisms into five kingdoms: animals, plants, fungi, protists, and
bacteria. The organizational scheme was based mainly on physical features, as opposed to physiology,
biochemistry, or molecular biology, all of which are used by modern systematics. The pioneering work of
American microbiologist Carl Woese in the early 1970s has shown, however, that life on Earth has evolved
along three lineages, now called domains—Bacteria, Archaea, and Eukarya. The first two are prokaryotic
cells with microbes that lack membrane-enclosed nuclei and organelles. The third domain contains the
eukaryotes and includes unicellular microorganisms together with the four original kingdoms (excluding
bacteria). Woese defined Archaea as a new domain, and this resulted in a new taxonomic tree (see this
figure). Many organisms belonging to the Archaea domain live under extreme conditions and are called
extremophiles. To construct his tree, Woese used genetic relationships rather than similarities based on
morphology (shape).
Woese’s tree was constructed from comparative sequencing of the genes that are universally distributed,
present in every organism, and conserved (meaning that these genes have remained essentially unchanged
throughout evolution). Woese’s approach was revolutionary because comparisons of physical features are
insufficient to differentiate between the prokaryotes that appear fairly similar in spite of their tremendous
biochemical diversity and genetic variability (Figure 1.18). The comparison of homologous DNA and RNA
sequences provided Woese with a sensitive device that revealed the extensive variability of prokaryotes,
and which justified the separation of the prokaryotes into two domains: bacteria and archaea."""

if (__name__ == "__main__"):
    import sys
    for elem in sys.argv:
        if (elem == "test_img"):
            gen_images(test_prompt, "test1", 0)
        elif (elem == "test_shorten_prompt"):
            shorten_prompt_iter(test_prompt, 50, 10, BedrockChat(model_id="meta.llama2-70b-chat-v1", client = bedrock_runtime))
        elif (elem == "test_prompt_integration"):
            print(chat_generate_script("biology", test_prompt))
        elif (elem == "test_full_integration"):
            uuid = str(uuid4())
            gen_video(gen_content(chat_generate_script("biology", test_prompt), uuid), uuid)

