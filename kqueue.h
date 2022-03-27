// @Author czkboy
// @Email czkboy000229@gmail.com
#pragma once
#define OS_MACOSX

#ifdef OS_LINUX
#include <sys/epoll.h>
#elif defined(OS_MACOSX)

#include <sys/event.h>
#else
#error "platform unsupported"
#endif
#include <poll.h>
#include <memory>
#include <unordered_map>
#include <vector>
// #include "Channel.h"
#include "HttpData.h"
#include "mactimer.h"

const int kReadEvent = POLLIN;
const int kWriteEvent = POLLOUT;

class PollerBase {
 public:
  PollerBase()= default;
  virtual ~PollerBase()= default;
  virtual void epoll_add(SP_Channel request, int timeout)=0;
  virtual void epoll_mod(SP_Channel request, int timeout)=0;
  virtual void epoll_del(SP_Channel request)=0;
  std::vector<std::shared_ptr<Channel>> poll();
  std::vector<std::shared_ptr<Channel>> getEventsRequest(int events_num);
  virtual void add_timer(std::shared_ptr<Channel> request_data, int timeout)=0;
  virtual int getEpollFd()=0;
  virtual void handleExpired()=0;

 private:
  static const int MAXFDS = 100000;
  int epollFd_;
  // std::vector<epoll_event> events_;
  std::shared_ptr<Channel> fd2chan_[MAXFDS];
  std::shared_ptr<HttpData> fd2http_[MAXFDS];
  TimerManager timerManager_;
};