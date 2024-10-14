# Comandos Despliegue en AWS

## Compilar Angular

```jsx
npm rum build
```

## Política S3

```jsx
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PublicReadGetObject",
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "s3:GetObject"
            ],
            "Resource": [
                "arn:aws:s3:::Bucket-Name/*"
            ]
        }
    ]
}
```

## Verificar si tiene git

```jsx
git --version
```

## Instalar GO

```jsx
sudo apt-get update
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -xvf go1.21.0.linux-amd64.tar.gz
sudo mv go /usr/local
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
source ~/.profile
go version
```

## GitHub

```jsx
git clone https://github.com/AndresPontaza/MIA_D_LABORATORIO_S2_2024.git
git pull
```

## Link del sitio web en AWS

http://mia-ejemplo9.s3-website.us-east-2.amazonaws.com