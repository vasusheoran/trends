#   Building the image
#       docker build -f node.dockerfile -t gokusayon/trends-dashboard . 
#   Running container after linking with existing container
#       docker run -p 80:80 --name dash -d --link trends:trends-app gokusayon/trends-dashboard    
#   Go to container shell
#       docker exec -i -t dash /bin/sh

#   Create an isolated network.
#       docker network create --driver bridge isolated_network 
#       docker run -d --net=isolated_network -p 5000:5000 --name trends gokusayon/trends-app
#       docker run -d --net=isolated_network -p 80:80 --name dash gokusayon/trends-dashboard

### STAGE 1: Setup ###
FROM        node:1.21.3-alpine

RUN         mkdir -p /ng-app
WORKDIR     /ng-app

COPY        package.json /ng-app

RUN         npm install -g @angular/cli @angular-devkit/build-angular && npm install

EXPOSE      4200
EXPOSE      49153

CMD         ["npm", "start"]