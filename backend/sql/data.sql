create database gpt_project;
create table gpt_project.users
(
    uid      int(11) unsigned auto_increment
        primary key,
    email    varchar(32) default ''                not null comment '账号邮箱',
    password varchar(32) default ''                not null comment '密码',
    token    varchar(32) default ''                not null comment '登录token',
    ct       datetime    default CURRENT_TIMESTAMP not null comment '创建时间',
    ut       datetime    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    constraint uniq_email
        unique (email)
)
    comment '用户表' charset = utf8;

create table gpt_project.experiment
(
    id     int(11) unsigned auto_increment
        primary key,
    topic  varchar(2048) default ''                not null comment '实验主题内容',
    status tinyint       default 0                 not null comment '实验状态,1为在线，0为下线',
    ct     datetime      default CURRENT_TIMESTAMP not null comment '创建时间',
    ut     datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间'
) comment '实验表' charset = utf8;

create table gpt_project.guide_page
(
    id                    int(11) unsigned auto_increment
        primary key,
    exid                  int           default 0                 not null comment '实验id',
    content               varchar(2048) default ''                not null comment '页面文本',
    next_button           tinyint       default 0                 not null comment '是否有下一页按钮，1为是0位否',
    yes_or_no             tinyint       default 0                 not null comment '是否有是否按钮，1为是0位否',
    chat                  tinyint       default 0                 not null comment '是否是聊天界面，1为是0位否',
    chat_times            int           default 0                 not null comment '限制聊天次数',
    chat_tips             varchar(100)  default ''                not null comment '聊天框提示标语',
    save_chat_info        tinyint       default 0                 not null comment '是否保存此页面的聊天信息,1为是，0为否',
    get_chat_from         int           default 0                 not null comment '从哪个页面id获取聊天信息',
    next_page             int           default 0                 not null comment '下一页id',
    yes_page              int           default 0                 not null comment '是按钮的id',
    no_page               int           default 0                 not null comment '否按钮的id',
    countdown             int           default 0                 not null comment '页面倒计时',
    answer_page           tinyint       default 0                 not null comment '是否是回答问题界面，1为是0位否',
    answer_time_countdown int           default 0                 not null comment '提交问题倒计时',
    ct                    datetime      default CURRENT_TIMESTAMP not null comment '创建时间',
    ut                    datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间'

) comment '引导信息表' charset = utf8;

create table gpt_project.questionnaire
(
    id    int(11) unsigned auto_increment
        primary key,
    exid  int           default 0                 not null comment '实验id',
    guide varchar(4096) default ''                not null comment '问卷引导介绍文本',
    ct    datetime      default CURRENT_TIMESTAMP not null comment '创建时间',
    ut    datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间'

) comment '问卷调查表' charset = utf8;
create table gpt_project.question
(
    id          int(11) unsigned auto_increment
        primary key,
    qnid        int           default 0                 not null comment '问卷调查id',
    type        int           default 0                 not null comment '题目类型，1为单选，2为多选，3为填空，4为评分题',
    is_required int           default 0                 not null comment '是否必做题，1为是，0为否',
    content     varchar(2048) default ''                not null comment '问题文本',
    choices     varchar(4096) default ''                not null comment '选择项集合，用;进行分隔',
    ct          datetime      default CURRENT_TIMESTAMP not null comment '创建时间',
    ut          datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间'

) comment '问卷问题表' charset = utf8;
create table gpt_project.question_score_id
(
    sid  int(11) unsigned auto_increment primary key,
    uid  int      default 0                 not null comment '用户uid',
    qnid int      default 0                 not null comment '问卷id',
    qid  int      default 0                 not null comment '问卷问题id',
    ct   datetime default CURRENT_TIMESTAMP not null comment '创建时间',
    ut   datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    UNIQUE KEY `uid_qid` (`uid`, `qid`)
) comment '评分问题id表' charset = utf8;

create table gpt_project.question_scoring
(
    id    int(11) unsigned auto_increment primary key,
    sid   int      default 0                 not null comment '评分题答案id',
    no    tinyint  default 0                 not null comment '评分题第几个选项',
    score tinyint  default 0                 not null comment '分数',
    uid   int      default 0                 not null comment '用户uid',
    qnid  int      default 0                 not null comment '问卷id',
    qid   int      default 0                 not null comment '问卷问题id',
    ct    datetime default CURRENT_TIMESTAMP not null comment '创建时间',
    ut    datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间'
) comment '评分题答案' charset = utf8;

create table gpt_project.experiment_reply
(
    id      int(11) unsigned auto_increment
        primary key,
    uid     int           default 0                 not null comment '用户uid',
    exid    int           default 0                 not null comment '实验id',
    content varchar(2048) default ''                not null comment '回答文本',
    access  tinyint       default 0                 not null comment '是否是访问时的记录，1为是0为否',
    ct      datetime      default CURRENT_TIMESTAMP not null comment '创建时间',
    ut      datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间'
) comment '实验回复表' charset = utf8;

create table gpt_project.question_reply
(
    id         int(11) unsigned auto_increment
        primary key,
    uid        int           default 0                 not null comment '用户uid',
    qnid       int           default 0                 not null comment '问卷id',
    qid        int           default 0                 not null comment '问卷问题id',
    reply_text varchar(2048) default ''                not null comment '填空题文本',
    choices    int           default 0                 not null comment '选择题答案选择,用二进制位代表实际选择,支持多选,如1代表选择了A，2代表选择了B，4代表选择了C，8代表选择了D，3=1+2代表选择了A+B',
    sid        int           default 0                 not null comment '评分题答案id',
    ct         datetime      default CURRENT_TIMESTAMP not null comment '创建时间',
    ut         datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间'
) comment '问卷问题选择填写表' charset = utf8;

create table gpt_project.experiment_chat_history
(
    id        int(11) unsigned auto_increment
        primary key,
    uid       int           default 0                 not null comment '用户uid',
    exid      int           default 0                 not null comment '实验id',
    query     varchar(4096) default ''                not null comment '用户问gpt的内容',
    gpt_reply varchar(4096) default ''                not null comment 'gpt回复内容',
    ct        datetime      default CURRENT_TIMESTAMP not null comment '创建时间',
    ut        datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间'
) comment 'gpt实验问答历史记录表' charset = utf8;

create table experiment_group
(
    id     int(11) unsigned auto_increment
        primary key,
    name   varchar(2048) default ''                not null comment '实验组名称',
    status tinyint       default 0                 not null comment '实验组状态,1为在线，0为下线',
    ct     datetime      default CURRENT_TIMESTAMP not null comment '创建时间',
    ut     datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间'
)
    comment '实验组表' charset = utf8;

create table group_member
(
    id        int(11) unsigned auto_increment
        primary key,
    group_id  int      default 0                 not null comment '实验组id',
    member_id int      default 0                 not null comment '实验组成员id',
    ct        datetime default CURRENT_TIMESTAMP not null comment '创建时间',
    ut        datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间'
)
    comment '实验组成员表' charset = utf8;