
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

## Uninstall charts
uninstall:
	$(HELM) uninstall $(NAME) -n $(NAMESPACE)

.PHONY : help install uninstall