INSERT INTO good_srv_brands
(create_time, update_time, deleted_at, is_deleted, name, logo)
VALUES
    (NOW(), NOW(), NULL, 0, '苹果', 'https://albelt-mxshop.oss-cn-shenzhen.aliyuncs.com/brand_logo/apple-14.svg'),
    (NOW(), NOW(), NULL, 0, '华为', 'https://albelt-mxshop.oss-cn-shenzhen.aliyuncs.com/brand_logo/huawei.svg'),
    (NOW(), NOW(), NULL, 0, '三星', 'https://albelt-mxshop.oss-cn-shenzhen.aliyuncs.com/brand_logo/samsung.svg'),
    (NOW(), NOW(), NULL, 0, '诺基亚', 'https://albelt-mxshop.oss-cn-shenzhen.aliyuncs.com/brand_logo/nokia-3.svg'),
    (NOW(), NOW(), NULL, 0, '联想', 'https://albelt-mxshop.oss-cn-shenzhen.aliyuncs.com/brand_logo/lenovo-logo-2015.svg'),
    (NOW(), NOW(), NULL, 0, '惠普', 'https://albelt-mxshop.oss-cn-shenzhen.aliyuncs.com/brand_logo/hp-5.svg');



INSERT INTO good_srv_categories
(create_time, update_time, deleted_at, is_deleted, name, parent_category_id, `level`, is_tab)
VALUES
    (NOW(), NOW(), NULL, 0, '数码', NULL, 1, 0);

INSERT INTO good_srv_categories
(create_time, update_time, deleted_at, is_deleted, name, parent_category_id, `level`, is_tab)
VALUES
    (NOW(), NOW(), NULL, 0, '手机', 1, 2, 0),
    (NOW(), NOW(), NULL, 0, '电脑', 1, 2, 0);

INSERT INTO good_srv_categories
(create_time, update_time, deleted_at, is_deleted, name, parent_category_id, `level`, is_tab)
VALUES
    (NOW(), NOW(), NULL, 0, '智能手机', 2, 3, 0),
    (NOW(), NOW(), NULL, 0, '老人机', 2, 3, 0),
    (NOW(), NOW(), NULL, 0, '笔记本电脑', 3, 3, 0),
    (NOW(), NOW(), NULL, 0, '台式电脑', 3, 3, 0);


INSERT INTO good_srv_category_brand_relations
(create_time, update_time, deleted_at, is_deleted, category_id, brands_id)
VALUES
    (NOW(), NOW(), NULL, 0, 4, 1),
    (NOW(), NOW(), NULL, 0, 4, 2),
    (NOW(), NOW(), NULL, 0, 4, 3),
    (NOW(), NOW(), NULL, 0, 5, 4),
    (NOW(), NOW(), NULL, 0, 6, 1),
    (NOW(), NOW(), NULL, 0, 6, 5),
    (NOW(), NOW(), NULL, 0, 6, 6),
    (NOW(), NOW(), NULL, 0, 7, 5),
    (NOW(), NOW(), NULL, 0, 7, 6);

INSERT INTO good_srv_banners
(create_time, update_time, deleted_at, is_deleted, image, url, `index`)
VALUES
    (NOW(), NOW(), NULL, 0,  'https://albelt-mxshop.oss-cn-shenzhen.aliyuncs.com/banners/1.jpg', '', 1),
    (NOW(), NOW(), NULL, 0,  'https://albelt-mxshop.oss-cn-shenzhen.aliyuncs.com/banners/2.jpg', '', 2),
    (NOW(), NOW(), NULL, 0,  'https://albelt-mxshop.oss-cn-shenzhen.aliyuncs.com/banners/3.jpg', '', 3);

INSERT INTO mxshop_srvs.good_srv_goods
(create_time, update_time, deleted_at, is_deleted, category_id, brand_id, on_sale, ship_free, is_new, is_hot, name, goods_sn, click_num, sold_num, fav_num, market_price, shop_price, goods_brief, images, desc_images, goods_front_image)
VALUES(NOW(), NOW(), NULL, 0, 4, 1, 1, 1, 0, 1, 'iphone13', 'sn-10001', 0, 0, 0, 6500, 5999, '苹果还是13香', '', '', '');

