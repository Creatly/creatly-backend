FROM nginx:1.15-alpine

COPY ./deploy/nginx/nginx.conf.template /
COPY ./deploy/nginx/docker-entrypoint.sh /

RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT ["sh", "/docker-entrypoint.sh"]

CMD ["nginx", "-g", "daemon off;"]