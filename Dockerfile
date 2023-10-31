## One Stage Build - Remove when 2-stage works
FROM amd64/ubuntu:latest

COPY institute-person-api /institute-person-api
COPY PATCH_LEVEL /PATCH_LEVEL

EXPOSE 8081
CMD ["/institute-person-api"]
