
create table category
(
    id          bigint unsigned                          not null comment '唯一标识'
        primary key,
    parent      bigint      default -1                   null comment '父级分类',
    name        varchar(20)                              not null comment '名称',
    level       tinyint                                  not null comment '级别',
    description varchar(255)                             not null comment '描述',
    create_time datetime(6) default CURRENT_TIMESTAMP(6) not null comment '创建时间',
    update_time datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6) comment '更新时间'
)
    comment '分类表';

create table tag
(
    id          bigint unsigned                          not null comment '唯一标识'
        primary key,
    name        varchar(20)                              not null comment '名称',
    create_time datetime(6) default CURRENT_TIMESTAMP(6) not null comment '创建时间',
    update_time datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6) comment '更新时间'
)
    comment '标签表';

create table book
(
    id          bigint unsigned                          not null comment '唯一标识'
        primary key,
    name        varchar(20)                              not null,
    author      varchar(20)                              not null,
    category    bigint                                   not null,
    create_time datetime(6) default CURRENT_TIMESTAMP(6) not null,
    update_time datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6)
)
    comment '书表';

create table book_tag
(
    id          bigint unsigned                          not null comment '唯一标识'
        primary key,
    book_id     bigint                                   not null comment '书',
    tag_id      bigint                                   not null,
    create_time datetime(6) default CURRENT_TIMESTAMP(6) not null comment '创建时间',
    update_time datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6) comment '更新时间'
)
    comment '标签表';

create table book_volume
(
    book_id     bigint unsigned                          not null comment '书ID',
    id          bigint unsigned                          not null comment '唯一标识'
        primary key,
    `index`     int unsigned                             not null comment '索引',
    name        varchar(20)                              not null comment '名称',
    create_time datetime(6) default CURRENT_TIMESTAMP(6) not null comment '创建时间',
    update_time datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6) comment '更新时间'
)
    comment '分卷表';

create table book_chapter
(
    book_id     bigint unsigned                          not null comment '书ID',
    volume_id   bigint unsigned                          not null comment '分卷ID',
    id          bigint unsigned                          not null comment '唯一标识'
        primary key,
    `index`     int unsigned                             not null comment '索引',
    name        varchar(20)                              not null comment '名称',
    create_time datetime(6) default CURRENT_TIMESTAMP(6) not null comment '创建时间',
    update_time datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6) comment '更新时间'
)
    comment '章节表';

create table book_paragraph
(
    book_id     bigint unsigned                          not null comment '书ID',
    volume_id   bigint unsigned                          not null comment '分卷ID',
    chapter_id  bigint unsigned                          not null comment '章节ID',
    id          bigint unsigned                          not null comment '唯一标识'
        primary key,
    `index`     int unsigned                             not null comment '索引',
    content     text                                     not null comment '段落内容',
    create_time datetime(6) default CURRENT_TIMESTAMP(6) not null comment '创建时间',
    update_time datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6) comment '更新时间'
)
    comment '段落表';

