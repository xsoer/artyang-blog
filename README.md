- go和vue的结合
- 包含两个vue,一个前端项目，一个后端项目
- nginx的配置
```
server {
    listen 80;
    server_name local.artyang.vip
    index index.html index.htm index.php;
    # root /Users/zxc/GOPRO/src/github.com/xsoer/artyang-blog;

    location /api {
        root /Users/zxc/GOPRO/src/github.com/xsoer/artyang-blog;
        # rewrite         ^.+api/?(.*)$ /$1 break;
        proxy_pass      http://127.0.0.1:8080;
    }
    location /x {
        root /Users/zxc/GOPRO/src/github.com/xsoer/artyang-blog;
        proxy_pass      http://127.0.0.1:8080;
    }

    # --------------------admin-------------------------
    location =/admin {
        root /Users/zxc/GOPRO/src/github.com/xsoer/artyang-blog/views-admin;
        #if ( !-f $document_root/index.html ) {
               proxy_pass              http://127.0.0.1:8091;
        #}
        proxy_redirect          off;
        proxy_set_header        Host            $host;
        proxy_set_header        X-Real-IP       $remote_addr;
        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 热加载
    location =/admin/__webpack_hmr {
        root /Users/zxc/GOPRO/src/github.com/xsoer/artyang-blog/views-admin;
        proxy_pass              http://127.0.0.1:8091;
        proxy_http_version      1.1;
        proxy_set_header        Upgrade $http_upgrade;
        proxy_set_header        Connection "Upgrade";
    }

    location ^~ /admin/static/ {
        root /Users/zxc/GOPRO/src/github.com/xsoer/artyang-blog/views-admin;
        #if ( !-d $document_root/static ) {
                proxy_pass              http://127.0.0.1:8091;
        #}
        proxy_redirect          off;
        proxy_set_header        Host            $host;
        proxy_set_header        X-Real-IP       $remote_addr;
        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # ------------------front---------------------------
    location =/ {
        root /Users/zxc/GOPRO/src/github.com/xsoer/artyang-blog/views-front;
        # if ( !-f $document_root/index.html ) {
               proxy_pass              http://127.0.0.1:8081;
        # }
        proxy_redirect          off;
        proxy_set_header        Host            $host;
        proxy_set_header        X-Real-IP       $remote_addr;
        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 热加载
    location =/__webpack_hmr {
        root /Users/zxc/GOPRO/src/github.com/xsoer/artyang-blog/views-front;
        proxy_pass              http://127.0.0.1:8081;
        proxy_http_version      1.1;
        proxy_set_header        Upgrade $http_upgrade;
        proxy_set_header        Connection "Upgrade";
    }

    location ^~ /static/ {
        root /Users/zxc/GOPRO/src/github.com/xsoer/artyang-blog/views-front;
        # if ( !-d $document_root/static ) {
                proxy_pass              http://127.0.0.1:8081;
        # }
        proxy_redirect          off;
        proxy_set_header        Host            $host;
        proxy_set_header        X-Real-IP       $remote_addr;
        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
    }

}
```