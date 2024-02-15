BASE=report
DVI=${BASE}.dvi
PS=${BASE}.ps
LOG=${BASE}.log
AUX=${BASE}.aux

.PHONY: clean

project3_pdf:
	./prefile1.sh
	./bin1.sh
	./run1.sh
	./report.sh

clean:
	rm -f ${DVI} ${PS} ${LOG} ${AUX}