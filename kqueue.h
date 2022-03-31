#pragma once
#include <assert.h>
#include <poll.h>
#include <string.h>
#include <sys/time.h>
#include <sys/types.h>
#include <atomic>
#include <map>
#include <zinx.h>
#include <base/noncopyable.h>
namespace zinx {

const int kMaxEvents = 2000;
const int kReadEvent = POLLIN;
const int kWriteEvent = POLLOUT;
struct Channel :public Ichannel{};

struct PollerBase : private noncopyable {
    int64_t id_;
    int lastActive_;
    PollerBase() : lastActive_(-1) {
        static std::atomic<int64_t> id(0);
        id_ = ++id;
    }
    virtual void addChannel(Channel *ch) = 0;
    virtual void removeChannel(Channel *ch) = 0;
    virtual void updateChannel(Channel *ch) = 0;
    virtual void loop_once(int waitMs) = 0;
    virtual ~PollerBase(){};
};

PollerBase *createPoller();

}  // namespace zinx