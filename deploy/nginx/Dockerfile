FROM nginx:1.15-alpine

COPY ./nginx.conf.template /
COPY ./docker-entrypoint.sh /

RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT ["sh", "/docker-entrypoint.sh"]

CMD ["nginx", "-g", "daemon off;"]