#   Building the image
#       docker build -f Dockerfile -t gokusayon/trends-dashboard . 
#   Running container after linking with existing container
#       sudo docker run -p 80:80 --name dash -d --link trends:trends-app gokusayon/trends-dashboard    
#       sudo docker run -p 80:80 --name dash -d gokusayon/trends-dashboard   
#   Go to container shell
#       docker exec -i -t dash /bin/sh
#   sudo docker build -f Dockerfile -t gokusayon/trends-app .
#   sudo docker run -d -p 5000:5000 --name trends gokusayon/trends-app
#   sudo docker run -d -p 80:4200    --name dash -v $(pwd):/var/www -w "/var/www" node npm start
#   sudo docker run -d -p 80:4200 --link trends:trends --name dash -v $(pwd):/var/www -w "/var/www" node npm start 
#   sudo docker run -d -p 80:4200 --network isolated_network --name dash -v $(pwd):/var/www -w "/var/www" node npm start

# ### STAGE 1: Build ###
# FROM        node:14-alpine3.12 as build-stage
# LABEL       AUTHOR="Vasu Sheoran"  

# ## Storing node modules on a separate layer will prevent unnecessary npm installs at each build
# RUN         mkdir /ng-app
# WORKDIR     /ng-app

# COPY        package.json /ng-app
# RUN         npm install @angular/cli@7.3.9 
# RUN         npm install

# ## Build the angular app in production mode and store the artifacts in dist folder
# COPY        . /ng-app
# RUN         npm run build:prod


### STAGE 2: Setup ###
FROM        nginx:1.21.3-alpine
RUN         rm -rf /usr/share/nginx/html/*

COPY        default.conf /etc/nginx/conf.d/default.conf
COPY        dist /usr/share/nginx/html
# COPY        --from=build-stage dist /usr/share/nginx/html

EXPOSE 4200 80


CMD         ["nginx", "-g", "daemon off;"]