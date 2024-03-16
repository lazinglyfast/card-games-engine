# Description

- A card games engine with an API that allows for manipulation of a standard deck of cards
- The API can be succintly described as:
  - GET `http://localhost/create?cards=A2,8C&shuffled` where cards and shuffled are optional
  - GET `http://localhost/open/{guid}`
  - GET `http://localhost/draw/{guid}?count=2` where count is optional and defaults to 1

# Running

- git clone `git@github.com:lazinglyfast/card-games-engine.git`
- in a terminal: cd server && go run . // this runs the server
- in another terminal: cd js-client-to-go-cards && npm run dev
- one can interact with the server with a client that can be:
  - via command line with for instance "curl http://localhost:8000/create"
  - with an app like postman
  - with the simple react app launched above by visiting: http://localhost:5173/
    ![alt text](https://github.com/lazinglyfast/card-games-engine/blob/main/react_app.png?raw=true)

# Comments on the evaluation

- extensibility/complexity
  - we could make a deck so extensible that it could work with any number of cards, suits and ranks or include other concepts entirely (i.e. a healing card) but if that's not an immediate or foreseeable requirement there's no need to over-engineer
  - complexity must be tamed and one of the most effective ways to do that is to not add more code
  - one point of extension might be the addition of joker cards. That would require us to rethink the whole design. In the pursue of simplicity no accomodations were made for that scenario
- comments
  - several comments have been added with the evaluation in mind like design decisions and coding style. In an actual production setting I tend to only add comments when I've failed to express intent using code only.
  - another use case for comments are documentation
  - yet another are sneaky edge cases like: "using the ASCII subset of UTF-8 so this is ok"
- testing
  - not familiar with go's bdd testing frameworks. Lookend into goconvey and ginkgo but they felt awkward to use so chose to leave it out
  - added regular unit tests for the deck package and http tests for the http server
  - if time had allowed I'd like to also have added some end-to-end testing treating our API as a blackbox
- data storage
  - for this evaluation I kept the state in memory but obviously in a real application state would be persisted to a database
- there are so many finer points that such an app should consider but they are obviously out-of-scope like
  - authentication so that one user cannot mess with another's deck
  - port number should not be hardcoded and should be dynamic
  - coming up with a design that includes joker cards or more esoteric cards and card features other than rank and suit
  - persistent storage
  - CI/CD
  - containerization
