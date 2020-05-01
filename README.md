# EXCommerce

Example Commerce System

## Development requirements

* Docker
* make

## API

* Start api doc server on [localhost:8081](http://localhost:8081/): `make redoc`

## Design

```
+--------------------------------------------+          +---------------------------------------------+
| SPA Frontend (Vue)                         |          |                        RESTful Backend (Go) |
|                                            |          |                                             |
|  +---------------------------+   +---------+--+    +--+--+  +---------------+--------------------+  |
|  | App Logic                 |   | API Client |    | API +--> Registration  | Controllers        |  |
|  |                           |   |            +---->     |  +---------------+                    |  |
|  |                           |   |            |    |     +--> Login         |                    |  |
|  |                           +--->            +---->     |  +---------------+                    |  |
|  |                           |   |            |    |     +--> Store Cart    |                    |  |
|  +---+-----------------------+   |            +---->     |  +---------------+                    |  |
|      |                           |            |    |     +--> ...           |                    |  |
|      |         +-----------------+ Auth       |    |     |  +---------------+                    |  |
|      |         |                 |            |    |     |                  |                    |  |
|      |         |          +------>            |    |     |   +------------+ |                    |  |
|      |         |          |      |            |    |     +---> Basic Auth | |                    |  |
|      |         |          |      +---------+--+    +--+--+   +----------+-+ +-------+-+-+--------+  |
|      |         |          |                |          |                 |           | | |           |
|  +---v---------v---+  +---+------------+   |          |  +--------------------------v-v-v--------+  |
|  | Persistence     |  | Server+side    |   |          |  | Models       |                        |  |
|  |                 |  |                |   |          |  |              |                        |  |
|  +-------+-+-+-----+  | +------+ +---+ |   |          |  |  +------+  +-v----+  +-------+  +---+ |  |
|          | | |        | | Cart | |...| |   |          |  |  | Cart |  | User |  | Order |  |...| |  |
|          | | +--------> +------+ +---+ |   |          |  |  +------+  +------+  +-------+  +---+ |  |
|          | |          |                |   |          |  |                                       |  |
|          | |          +----------------+   |          |  +---------------+-+-+-------------------+  |
|          | |                               |          |                  | | |                      |
|          | +-------------------+           |          |  +---------------v-v-v-------------------+  |
|          |                     |           |          |  | Persistence                           |  |
|  +-------v---------+  +--------v-------+   |          |  |                                       |  |
|  | Session Storage |  | Local Storage  |   |          |  +-----+-------------+-------------+-----+  |
|  |                 |  |                |   |          |        |             |             |        |
|  | +-------------+ |  | +------+ +---+ |   |          |  +-----v----+  +-----v-----+  +----v-----+  |
|  | | Credentials | |  | | Cart | |...| |   |          |  | Database |  | In-memory |  | ...      |  |
|  | +-------------+ |  | +------+ +---+ |   |          |  |          |  |           |  |          |  |
|  |                 |  |                |   |          |  |          |  |           |  |          |  |
|  +-----------------+  +----------------+   |          |  +----------+  +-----------+  +----------+  |
|                                            |          |                                             |
+--------------------------------------------+          +---------------------------------------------+
```
