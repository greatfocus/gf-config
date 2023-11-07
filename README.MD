# [Install Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux)
Download the latest release with the command

    curl -LO https://dl.k8s.io/release/v1.28.3/bin/linux/amd64/kubectl

Validate the binary (optional)

    echo "$(<kubectl.sha256)  kubectl" | sha256sum --check

Install kubectl

    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

Test to ensure the version you installed is up-to-date

    kubectl version --client


# [Install Minikube](https://minikube.sigs.k8s.io/docs/start)
Installation
    
    curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 sudo install minikube-linux-amd64 /usr/local/bin/minikube

Enable ingress

    minikube addons enable ingress

Start minikube

    minikube start

# GitFlow Actions
Pull Request

    - First time PR creation runs pull_request_builder.yaml
    - Pushing a commit to an existing PR runs pull_request_builder.yaml
    - Merging an existing PR to develop/master runs continuous_deployment.yaml
    - Direct commit push to develop/master runs continuous_deployment.yaml

Release/Hotfix

    - First time release/hotfix creation against develop/master runs pull_request_builder.yaml
    - Pushing a commit to an existing release/hotfix runs pull_request_builder.yaml
    - Creating a new PR against the release/hotfix runs pull_request_builder.yaml
    - Pushing a commit to an existing PR that tracks a release/hotfix runs pull_request_builder.yaml
    - Merging an existing PR to a release/hotfix runs pull_request_builder.yaml
    - Direct commit push to release/hotfix runs pull_request_builder.yaml
    - Merging an existing release/hotfix to develop/master runs continuous_deployment.yaml
    - Merging an existing release/hotfix to develop/master in CLI and pushing develop/master runs continuous_deployment.yaml

Tag

    - Manually creating a tag does nt trigger any pipeline
    - Pushing a tag from local does not trigger any pipeline
    
Manual Deployment

    - Manually deploying a specific git ref (branch/tag/sha) to a specific environment

Others

    - Closed PRs will not trigger any pipeline
    - Rollback is same as manual deployment so deploy a previous tag

# Add Github Secrets and token

Add Docker secrets and Kube config to GitHub. This will let your workflow to push images to DockerHub registry.

    DOCKERHUB_USERNAME
    DOCKERHUB_TOKEN


# Set minikube K8S cluster access to Github
This will let your workflows to deploy the application to your local K8S cluster.

Enable routing public requests to K8S API.

    kubectl proxy --port=8001 --accept-hosts='.*\.ngrok.io'

Expose kubectl proxy to the Internet with ngrok.
     
    ./ngrok http 127.0.0.1:8001
    Forwarding                    http://c309b9b1359f.ngrok.io -> http://127.0.0.1:8001
    Forwarding                    https://c309b9b1359f.ngrok.io -> http://127.0.0.1:8001

Create a modified copy of your local kube config.
     
    kubectl config view --flatten > ~/Desktop/kube_config
    
    Remove `certificate-authority-data` line
    Add `insecure-skip-tls-verify: true` line

    Replace `server` value to `https://c309b9b1359f.ngrok.io`
     
Finally it should look like something like below.
    

    apiVersion: v1
    clusters:
    - cluster:
        insecure-skip-tls-verify: true
        server: https://c309b9b1359f.ngrok.io
      name: minikube
    contexts:
    - context:
        cluster: minikube
        user: minikube
      name: minikube
    current-context: minikube
    kind: Config
    preferences: {}
    users:
    - name: minikube
      user:
        client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0.......
        client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0.......

     
Copy this content into `CLUSTER_CONFIG` GitHub secret. Also additional ones for the other environments.


# killing sturbborn volumeClaim

    kubectl patch pvc PVC_NAME -p '{"metadata":{"finalizers": []}}' --type=merge

# Killing sturbborn ingress

    https://devopscube.com/setup-ingress-kubernetes-nginx-controller/

    kubectl get ingress --all-namespaces
    
    kubectl delete all  --all -n qa

# Generate Self Signed Certificate

    Step 1: Generate a CA private key

        $ openssl genrsa -out ca.key 2048

    Step 2: Create a self-signed certificate, valid for 365 days.

        $ openssl req -x509 \
            -new -nodes  \
            -days 365 \
            -key ca.key \
            -out ca.crt \
            -subj "/CN=*.localhost.com" \
            -addext "subjectAltName = DNS:*.localhost.com" 

    openssl base64 -in ca.crt -out ca.txt


# Generate Self Signed Certificate

Add the ingress host to host computer 

    sudo cat /etc/hosts

