# mirango-go #

mirango-go is a small and simple content management system for blogs built with Go, [chi](https://github.com/go-chi/chi), and [pq](https://github.com/lib/pq).

This repository houses the backend of mirango.io (RESTful API)

PostgreSQL > 9.6 is required with migrations handled by [mattes/migrate](https://github.com/mattes/migrate)

## How do I get set up? ##

* Ensure PostgreSQL is installed and running
* Set environmental variables
* Install third party go packages
* Initialize database with `mattes/migrate`

## Database Design ##
```sql
CREATE TYPE user_role AS ENUM ('ADMIN', 'MEMBER');

CREATE TABLE users (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  uname          text NOT NULL,
  role           user_role NOT NULL DEFAULT 'MEMBER',
  digest         text NOT NULL,
  email          text NOT NULL,
  gpg_key        text NOT NULL,
  last_online_at timestamptz NOT NULL DEFAULT NOW(),
  created_at     timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE posts (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        uuid NOT NULL REFERENCES users(id),
  title          text NOT NULL,
  slug           text NOT NULL,
  sub_title      text NOT NULL,
  short          text NOT NULL,
  post_content   text NOT NULL,
  digest         text NOT NULL,
  published      boolean DEFAULT FALSE,
  updated_at     timestamptz NOT NULL DEFAULT NOW(),
  created_at     timestamptz NOT NULL DEFAULT NOW()
)

CREATE TABLE images (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        uuid NOT NULL REFERENCES users(id),
  url            text NOT NULL,
  medium         text NOT NULL,
  small          text NOT NULL,
  caption        text NOT NULL,
  updated_at     timestamptz NOT NULL DEFAULT NOW(),
  created_at     timestamptz NOT NULL DEFAULT NOW()
)

CREATE TABLE tags (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name          text NOT NULL,
  slug          text NOT NULL
)

CREATE TABLE posts_tags (
  post_id       uuid NOT NULL REFERENCES posts(id),
  tag_id        uuid NOT NULL REFERENCES tags(id)
)
```

## Who do I talk to? ##

* Contact Sean @ [github.com/sdwalsh](https://www.github.com/sdwalsh) or [bitbucket.org/sdwalsh](https://www.bitbucket.org/sdwalsh)

## Contributing ##

Feel free to fork this repository and send pull requests!

# License #

MIT License

Copyright (c) 2017 Sean Walsh

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
