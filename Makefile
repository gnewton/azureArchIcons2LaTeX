
all:
	cd azure2tex; go build
	azure2tex/azure2tex
	cd tex; make; mv azureAllIcons.pdf ..

clean:
	cd tex; make clean
	rm *~
