# L.A.T.E.

The Language Agnostic Template Engine is provided as a shared library with C bindings and a command line tool.

It has many similarities to Shopify's Liquid templating language. After many years of using Liquid, I'm building LATE to keep many of Liquids benefits while fixing a number of it's shortcomings.

Also, I want to provide this to people outside of the Ruby ecosystem.

Goals of LATE

* Strict and well defined scoping rules
* Strict and well defined parsing rules
* Fully configurable filters and tags
* Multiple levels of error handling (super strict where any error will halt rendering, or try to render and just warn on issues)
* Great error reporting with the ability to try to render something no matter what
* Fast
* File format agnostic, but supports HTML-style escaping
