# ExCommerce

Example Commerce System

## Development requirements

* Docker
* make

## API

* Build the server: `make build`
* Start the server build on `localhost:8080`: `make run`
* Start api doc server on [localhost:8081](http://localhost:8081/): `make redoc`

### Configuration via environment variables

* `COUPON_DEFAULT_LIFETIME`: The default lifetime when creating a coupon and no
  expires at date is given. Use values like `10s`, `2.5m` or `1h30m` to express
  a duration. Defaults to `10s`.

## Frontend

* Run the dev frontend on [localhost:3000](http://localhost:3000/): `make frontend`

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
|  | Persistence     |  | Server-side    |   |          |  | Models       |                        |  |
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

Guest orders a product:

```
+----------------------+     +----------------------+     +----------------------+
| User                 |     | Frontend             |     | Backend              |
+----------------------+     +----------------------+     +----------------------+
|                      |     |                      |     |                      |
| adds product to cart +-----> stores cart in local |     |                      |
|                      |     | storage              |     |                      |
|                      |     |                      |     |                      |
| goes to checkout     +-----> requires login       |     |                      |
|                      |     |                      |     |                      |
|                      <-----+ show login           |     |                      |
|                      |     |                      |     |                      |
| want to create acc   +----->                      |     |                      |
|                      |     |                      |     |                      |
|                      <-----+ show register        |     |                      |
|                      |     |                      |     |                      |
| register             +----------------------------------> create user          |
|                      |     |                      |     |                      |
|                      |     | store user creds     <-----+                      |
|                      |     |                      |     |                      |
|                      |     | store cart on server +-----> apply discounts &    |
|                      |     |                      |     | store cart           |
|                      |     |                      |     |                      |
|                      |     | update local cart    <-----+                      |
|                      |     |                      |     |                      |
|                      <-----+ show checkout        |     |                      |
|                      |     |                      |     |                      |
| enter coupon         +------- add cart info ------------> validate order       |
|                      |     |                      |     |                      |
|                      <-----+ show order details   <-----+ order details        |
|                      |     |                      |     |                      |
| submit payment info  +------- add cart info ------------> lock cart, validate  |
|                      |     |                      |     | and store order      |
|                      |     |                      |     |                      |
|                      |     | clear cart           <-----+ order details        |
|                      |     |                      |     |                      |
|                      <-----+ show order details   |     |                      |
|                      |     |                      |     |                      |
+----------------------+     +----------------------+     +----------------------+
```

## Attribution

Images used:

* Orange by [Marcos Ramírez](https://unsplash.com/@marcosramirez_x?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText) on [Unsplash](https://unsplash.com/s/photos/orange?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText)
* Pear by [Amirul Islam](https://unsplash.com/@picturedesign96?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText) on [Unsplash](https://unsplash.com/s/photos/pear?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText)
* Apple by [Louis Hansel @shotsoflouis](https://unsplash.com/@louishansel?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText) on [Unsplash](https://unsplash.com/s/photos/apple?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText)
* Banana by [Mike Dorner](https://unsplash.com/@dorner?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText) on [Unsplash](https://unsplash.com/s/photos/banana?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText)
