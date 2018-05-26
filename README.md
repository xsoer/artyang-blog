- go和vue的结合
- nginx的配置
```
server {
    listen 80;
    server_name local.artyang.vip
    index index.html index.htm index.php;
    root /Users/zxc/GOPRO/src/github.com/xsoer/artyang-blog;

    location =/ {
        if ( !-f $document_root/index.html ) {
               proxy_pass              http://127.0.0.1:8081;
        }
        proxy_redirect          off;
        proxy_set_header        Host            $host;
        proxy_set_header        X-Real-IP       $remote_addr;
        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 热加载
    location =/__webpack_hmr {
        proxy_pass              http://127.0.0.1:8081;
        proxy_http_version      1.1;
        proxy_set_header        Upgrade $http_upgrade;
        proxy_set_header        Connection "Upgrade";
    }

    location ^~ /static/ {
        if ( !-d $document_root/static ) {
                proxy_pass              http://127.0.0.1:8081;
        }
        proxy_redirect          off;
        proxy_set_header        Host            $host;
        proxy_set_header        X-Real-IP       $remote_addr;
        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /api {
        rewrite         ^.+api/?(.*)$ /$1 break;
        proxy_pass      http://127.0.0.1:8080;
    }

}
```