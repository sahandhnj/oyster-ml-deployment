# Oyster(**) 

Oyster is a tool to train, deploy as API and manage your ML models. Each ML model has its own isolated life-cycle which is fully customable.

## Getting Started

In order to start an oyster project move to a parent directory of your model and execute the oyster init command. 
```
* oyster init –modelPath kerasmodel –name sentiment
```
modelPath is the parent directory of the trained model and preprocessing functions.
```
* ls kerasmodel
* model.json  model_weights.h5  preprocessing.py
````

Below you can see the default file structure:
```
* oyster init –modelPath keraspipeline –name sentiment
```

After iniating an oyster project the configuration of the model and it's versions will be stored in .oyster.

### Prerequisites
Docker
Kubernetes

## Deployment
Trained model can be deployed into an API either localy in docker containers or on kubernetes cluster over multiple networks.

### resource
resources are physical or containarized servers allocated for training and serving. Local or cloud kuberentes clusters can be added as resources.


### Version
In order to deploy a model a version needs to be commited. Each version has its own state and can be deployed into different stages (train, test, api, batch).


```
* oyster version commit 
* sentiment:v1 has been created. 
*
* oyster deploy --version 1 --resource local-docker --stage api
* sentiment:v1 has been deployed.
*
* predict: /dep/sentiment/v1/predict
* API traffic : /monitor/sentiment/v1/traffic?from=&to= 
*
* In order to setup auhtentication run command with --auth

```

### Cloud
In order to deploy models to cloud you can add kubernetes nodes as resources. It is also possible to configure oyster to scale up by creating more pods and nodes.

```
* oyster authenticate gcloud --authFile gcloud.yml
* oyster k8 deploy -c glcoud 
* 8 new resources has been added.
```

## Versioning
Currently oyster is in alpha version and is being tested and developed everyday. Regular updates are being released.

## Authors

* **Sahand Hosseininejad**
* **Bram Bloks**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

