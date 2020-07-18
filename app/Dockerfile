#   Building the image
#       sudo docker build -f Dockerfile -t gokusayon/trends-app .
#   Running container after linking with existing container
#       sudo docker run -p 5000:5000 -d --name trends gokusayon/trends-app 
#       sudo docker run -d --name trends gokusayon/trends-app 
#       sudo docker run -d --network isolated_network --name trends gokusayon/trends-app

FROM python:3.8-slim

LABEL AUTHOR="Vasu Sheoran"  

COPY ./req.txt /home/app/
RUN pip install -r ./home/app/req.txt

COPY . /home/app/
WORKDIR /home/app/

# ENV TZ=Asia/Kolkata
# RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
# EXPOSE 5000
CMD [ "python", "main.py" ]