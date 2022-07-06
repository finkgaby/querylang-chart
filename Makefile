
NAME				:= querylang-chart
NAMESPACE			:= default
HELM                := helm

help:
	@echo "install    : To install the querylang-chart subchart"
	@echo "uninstall  : To uninstall the querylang-chart subchart"

build:
	nerdctl build --namespace k8s.io  -t querylang:v1.1 .

## Installs Charts to kubernetes cluster
install:
	$(HELM) install --create-namespace -n $(NAMESPACE) $(NAME) .

all: build install

rmi:
	nerdctl rmi -f --namespace k8s.io querylang:v1.1
## Uninstall charts
uninstall:
	$(HELM) uninstall $(NAME) -n $(NAMESPACE)

all-remove: uninstall rmi


.PHONY : help install uninstall