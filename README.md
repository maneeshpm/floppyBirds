# FloppyShawtty 
A Flappy birds profile for the legendary J-Shaw

![Gameplay](/res/imgs/gameplay/img2.png?raw=true)
![Gameplay](/res/imgs/gameplay/img1.png?raw=true)
![Gameplay](/res/imgs/gameplay/img3.png?raw=true)

## Installing and Running

- Create Executable with `go build`, create an `app`
- Run the app with `./app`

<i> Tip combine both `go build && ./app` for the impatient </i>

## Contributing guidelines

New profiles can be created by simply creating new images/frames and pasting them inside `res` directory
The game executes in the `main.go` file. Game handler and physics are defined in `scene.go`. Other files `bird.go`, `pipe.go`, `card.go` are simply models for the respective entity. 




