mindis -- A Decidedly Minimal Redis Library for Go
==================================================

This is a minimalist redis implementation that is intended to provide a stable API that is orthogonal to other Go         packages.  It is better defined by listing the things it deliberately does not implement:

* no connection pools -- complex and hard to predict in concurrent applications
* no per-command methods -- a few simple API's support everything in redis
* no complicated pipelines -- since Mindis supports sending several commands at once before pausing to examine results

It was written after subtle changes in how Godis handled commands caused regressions and it was determined that Godis had a more experimental bent.  Despite this philosophical difference, Mindis is comparable in speed to Godis.


The Rationale of Send(), Next(), Scan():
========================================

It is difficult in a generic database api to predict all usage models. Redis, in particular, presents a challenge due to  its high performance usage model and diverse set of commands.  To provide support for the widest range of usage models,   Mindis uses a simple API where sending a command, waiting for a result, and decoding a successful result are three        distinct operations, which can each return an error.

Since the strategy of handling an error in any of these three steps is often the same, another method, Status() offers a way to access the first error that occurred after the last Send() call.  The error is also propagated forward, preventing Next and/or Scan methods from acting.

    var val int
    conn.Send("GET", "user_ct")
    conn.Next()
    err = conn.Scan(&val)

Note that this pattern can be more simply expressed using the convenient Exec() call.

    var val int
    err = conn.Exec(&val, "GET", "user_ct")



Using SUBSCRIBE and MONITOR:
============================

Certain redis commands, like SUBSCRIBE and MONITOR monopolize a redis connection indefinitely after use.  The pattern for using commands like this is to provide them with a dedicated connection and use a for loop of Next() and Scan() calls     until the connection closes or the sun burns out.
    
    t, err := net.Dial("tcp", "localhost:6379")
    if err != nil {
        println(err)
        return
    }
    defer t.Close()
    c := mindis.Wrap(c)
    c.Send("SUBSCRIBE", "events")
    m.Next() // eat the REDIS acknowledgement
    events := make(chan string)
    go process(events)
    for c.Next() == nil {
        var evt []string
        if c.Scan(&evt) != nil {
            break
        }
        events <- evt[2] // "message", "events", $message
    }
    println(m.Status().Error())


Using Deadlines:
================

Since Mindis simply wraps an existing io.ReadWriter, it is possible to use net.Conn instead, which implies the            availability of using SetWriteDeadline and SetReadDeadline to control how long Mindis will wait to perform an operation.

It is important to note that if a deadline is hit, the Conn state will no longer be consistent for subsequent reuse.  It  should be dropped and a new io.ReadWriter used.


Encoding and Decoding Values:
=============================

Mindis provides an explicit Marshaller and Unmarshaller interface to convert byte-arrays to and from values.  For values  that do not provide this interface, the encoding/json Marshaller and Unmarshaller defaults are tried, followed by support for basic string, bool, int and float conversions.

Notably, compared to Godis and other Redis libraries, Mindis does not default to using the "fmt" package; these encodings are not intended to be symmetric, and in the case of lists and maps, is almost always the wrong result.  A TYPE_ERROR is  produced in these cases.
