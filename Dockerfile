FROM scratch
MAINTAINER Cheradenine Zakalwe <zdrakenot@gmail.com>

ADD gramarr /

ENTRYPOINT ["/gramarr"]
CMD ["--config=/config/config.json"]
