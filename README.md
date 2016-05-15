# Deapexer

![Logo](images/logo.png)

## Written by PyratLabs 2016, [http://pyrat.io/](http://pyrat.io/)

### License

Deapexer has been released under the MIT License

    Copyright (c) 2016 PyratLabs (http://pyrat.io)
    
    Permission is hereby granted, free of charge, to any person obtaining a 
    copy of this software and associated documentation files (the "Software"), 
    to deal in the Software without restriction, including without limitation 
    the rights to use, copy, modify, merge, publish, distribute, sublicense, 
    and/or sell copies of the Software, and to permit persons to whom the 
    Software is furnished to do so, subject to the following conditions:
    
    The above copyright notice and this permission notice shall be included in 
    all copies or substantial portions of the Software.
    
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR 
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, 
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER 
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING 
    FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER 
    DEALINGS IN THE SOFTWARE.

### Description

This application has a very specific use case for dealing with redirecting the 
apex record of a domain to a subdomain. The main reason for this is that a 
domain apex record must not be a CNAME, it must be an A record.

A good example of when to use this is when working with AWS Elastic Load 
Balancers when you want to avoid using Route 53. This has been something that 
I have had to work around with customers. Often the solution is to 301 redirect
in Apache on one of the load balanced EC2 instances. This causes 
disproportionate load on EC2s.

What DeApexer allows you to do is host one shared service that generates these 
redirects for any domain with a low cost to disk I/O. Configuration is stored 
in memory so this is fairly light on resources.

### Installation

#### Native Executable

Included is a Makefile so that you can build the binary and install it to the
default location. To do this, run:

    make
    sudo make install

The executable will be moved to /usr/local/bin/deapexer and the configuration
moved to /etc/deapexer/config.json

#### Dockerfile

You can also run deapexer in a docker container. The best practice for this is:

    docker build -t deapexer .
    docker run -d -p 80:80 --read-only --name "go-deapexer" deapexer
