create table book
(
    id          bigint unsigned                          not null comment '唯一标识'
        primary key,
    origin_id   varchar(20)                              not null,
    name        varchar(50)                              not null,
    author      varchar(20)                              not null,
    category    bigint                                   not null,
    create_time datetime(6) default CURRENT_TIMESTAMP(6) not null,
    update_time datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6)
)
    comment '书表';

create index author_name_index
    on book (author, name);

create index origin_id_index
    on book (origin_id);

create table book_chapter
(
    book_id      bigint unsigned                          not null comment '书ID',
    volume_index int unsigned                             not null comment '分卷索引',
    id           bigint unsigned                          not null comment '唯一标识'
        primary key,
    origin_id    varchar(20)                              not null,
    `index`      int unsigned                             not null comment '索引',
    name         varchar(50)                              not null comment '名称',
    create_time  datetime(6) default CURRENT_TIMESTAMP(6) not null comment '创建时间',
    update_time  datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6) comment '更新时间'
)
    comment '章节表';

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
    name        varchar(50)                              not null comment '名称',
    create_time datetime(6) default CURRENT_TIMESTAMP(6) not null comment '创建时间',
    update_time datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6) comment '更新时间'
)
    comment '分卷表';

create table book_word
(
    book_id         bigint unsigned                          not null comment '书ID',
    volume_index    int unsigned                             not null comment '分卷索引',
    chapter_index   int unsigned                             not null comment '章节索引',
    id              bigint unsigned                          not null comment '唯一标识'
        primary key,
    `index`         int unsigned                             not null comment '索引',
    word_book_index int unsigned                             not null comment '字符在书中的索引位',
    word            char                                     not null comment '字符',
    create_time     datetime(6) default CURRENT_TIMESTAMP(6) not null comment '创建时间',
    update_time     datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6) comment '更新时间'
)
    comment '字符表';

create index index_book_idx
    on book_word (`index`, book_id);

create table category
(
    id          bigint unsigned                          not null comment '唯一标识'
        primary key,
    parent      bigint      default -1                   null comment '父级分类',
    name        varchar(50)                              not null comment '名称',
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
    name        varchar(50)                              not null comment '名称',
    create_time datetime(6) default CURRENT_TIMESTAMP(6) not null comment '创建时间',
    update_time datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6) comment '更新时间'
)
    comment '标签表';

