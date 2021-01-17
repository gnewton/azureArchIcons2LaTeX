
all:
	cd azure2tex; go build
	-rm icons_tex/*	
	azure2tex/azure2tex
	cd tex; make; mv azureAllIcons.pdf ..
	cd examples; make

clean:
	cd tex; make clean
	rm *~
