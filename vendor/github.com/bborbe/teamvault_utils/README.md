# Teamvault Utils

## Generate config directory with Teamvault secrets

Install:

```
go get github.com/bborbe/teamvault_utils/bin/teamvault_config_dir_generator
```

Config:

```
{
    "url": "https://teamvault.example.com",
    "user": "my-user",
    "pass": "my-pass"
}
```

Run:

```
teamvault_config_dir_generator \
-teamvault-config="~/.teamvault.json" \
-source-dir=templates \
-target-dir=results \
-logtostderr \
-v=2
```

## Parse variable Teamvault secrets

Install:

```
go get github.com/bborbe/teamvault_utils/bin/teamvault_config_parser
```

Sample config:

```
foo=bar
username={{ "vLVLbm" | teamvaultUser }}
password={{ "vLVLbm" | teamvaultPassword }}
bar=foo 
```

Run:

```
cat my.config | teamvault_config_parser
-teamvault-config="~/.teamvault.json" \
-logtostderr \
-v=2
```

## Teamvault Get Username

Install:

```
go get github.com/bborbe/teamvault_utils/bin/teamvault_username
```

Run:

```
teamvault_username \
--teamvault-config ~/.teamvault-sm.json \
--teamvault-key vLVLbm
```

## Teamvault Get Password

Install:

```
go get github.com/bborbe/teamvault_utils/bin/teamvault_password
```

Run:

```
teamvault_password \
--teamvault-config ~/.teamvault-sm.json \
--teamvault-key vLVLbm
```

## Teamvault Get Url

Install:

```
go get github.com/bborbe/teamvault_utils/bin/teamvault_url
```

Run:

```
teamvault_url \
--teamvault-config ~/.teamvault-sm.json \
--teamvault-key vLVLbm
```

## Continuous integration

[Jenkins](https://jenkins.benjamin-borbe.de/job/Go-Teamvault-Utils/)

## Copyright and license

    Copyright (c) 2016, Benjamin Borbe <bborbe@rocketnews.de>
    All rights reserved.
    
    Redistribution and use in source and binary forms, with or without
    modification, are permitted provided that the following conditions are
    met:
    
       * Redistributions of source code must retain the above copyright
         notice, this list of conditions and the following disclaimer.
       * Redistributions in binary form must reproduce the above
         copyright notice, this list of conditions and the following
         disclaimer in the documentation and/or other materials provided
         with the distribution.

    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
    "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
    LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
    A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
    OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
    SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
    LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
    DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
    THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
    (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
    OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
