-- MySQL dump 10.13  Distrib 8.0.42, for Linux (x86_64)
--
-- Host: localhost    Database: DeliveryAppDB
-- ------------------------------------------------------
-- Server version	8.0.42

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `Images`
--

DROP TABLE IF EXISTS `Images`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Images` (
  `id` int NOT NULL AUTO_INCREMENT,
  `url` varchar(500) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `public_id` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=73 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Images`
--

LOCK TABLES `Images` WRITE;
/*!40000 ALTER TABLE `Images` DISABLE KEYS */;
INSERT INTO `Images` VALUES (28,'https://www.coca-cola.com/content/dam/brands/us/coca-cola/coca-cola-logo.png','2025-09-24 09:46:55',''),(29,'https://upload.wikimedia.org/wikipedia/commons/2/27/Coca_Cola_Flasche_-_Original_Taste.jpg','2025-09-24 09:46:55',''),(30,'https://product.hstatic.net/200000534989/product/dsc08341-enhanced-nr_1_e6d5d0a13c8f42c2bd7cea59e03ce199_master.jpg','2025-09-24 09:46:55',''),(31,'https://www.pepsi.com/en-us/uploads/images/twil-can.png','2025-09-24 09:46:55',''),(32,'https://www.oreo.com/images/hero/oreo-original.png','2025-09-24 09:46:55',''),(33,'https://cdn.lottemart.vn/media/catalog/product/cache/9b5f86ccf0bb6d794da3fb554015eb8c/s/t/sting_dau_330ml.jpg','2025-09-24 09:46:55',''),(34,'https://cdn.thtrue.vn/wp-content/uploads/2022/04/tra-xanh-0-do.jpg','2025-09-24 09:46:55',''),(35,'https://www.oreo.com/images/hero/oreo-original.png','2025-09-24 09:46:55',''),(36,'https://cdn.lottemart.vn/media/catalog/product/cache/9b5f86ccf0bb6d794da3fb554015eb8c/s/t/sting_dau_330ml.jpg','2025-09-24 09:46:55',''),(37,'https://bavifoods.com/thumbs/740x740x1/upload/product/cam-ep-5018.jpg','2025-09-24 09:46:55',''),(38,'https://bavifoods.com/thumbs/740x740x1/upload/product/cam-ep-5018.jpg','2025-09-24 09:46:55',''),(39,'https://cdn-i.vtcnews.vn/files/news/2019/01/22/-145625.jpg','2025-09-24 09:46:55',''),(43,'Chi-em-thi-nhau-lung-mua-dua-xiem-ve-uong-sau-tiem-phong-cua-hang-moi-ngay-ban-5000-qua-1-1631524190-570-width650height431.jpg (650×431)','2025-09-24 09:46:55',''),(44,'medium_20200513_094458_574364_nuoc_dua_max_1800x1800_jpg_095dc5e7ad.jpg (750×563)','2025-09-24 09:46:55',''),(45,'coconut-water-benefits-17218412875751213756362.jpg (800×562)','2025-09-24 09:46:55',''),(46,'https://file.hstatic.net/1000199715/file/uong-sua-sau-sinh-1_90f6b928e6084e7e87c4e7a89e1b1be3_grande.jpg','2025-09-24 09:46:55',''),(47,'khi_nao_nen_cho_be_uong_sua_1_4401cf044a.jpg (800×600)','2025-09-24 09:46:55',''),(48,'https://suckhoedoisong.qltns.mediacdn.vn/324455921873985536/2025/2/21/dau-nanh-1-1740125251401155246723.jpg','2025-09-24 09:46:55',''),(49,'glass-soy-milk_20dc83bb32164c49bd11a7d7b60b717b_grande.jpg (600×377)','2025-09-24 09:46:55',''),(50,'may-lam-sua-dau-nanh-1-1412734006024.jpg (500×455)','2025-09-24 09:46:55',''),(51,'https://baothainguyen.vn/file/e7837c027f6ecd14017ffa4e5f2a0e34/032023/1-boba-tea-recipe-using-fresh-tapioca-pearls-1024x1024-1677809524112848165864_20230305161118.jpeg','2025-09-24 09:46:55',''),(52,'https://www.cet.edu.vn/wp-content/uploads/2018/04/tra-sua-tu-lam.jpg','2025-09-24 09:46:55',''),(53,'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRXtYkBiYM_AEWth56eq5VeEwFGlnh_cYm7cw&s','2025-09-28 18:34:00',''),(54,'https://billballcoffeetea.com/upload/product/img3589-3619-8202.jpg','2025-09-28 18:34:00',''),(55,'https://vcdn1-suckhoe.vnecdn.net/2023/02/01/iced-coffee-table-jpeg-1675223-7169-5352-1675223880.jpg?w=1200&h=0&q=100&dpr=1&fit=crop&s=ExRy7hbEHS2p2f2oOgWCjA','2025-09-28 18:34:00',''),(56,'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQie4Hw_GApMLyRn34hNhlXh46_33_56ZcfMA&s','2025-09-28 18:34:00',''),(57,'https://cdn.tgdd.vn/2020/07/CookProductThumb/59-620x620-3.jpg','2025-09-28 18:34:00',''),(58,'https://product.hstatic.net/200000848723/product/combo_2_burger_e7394736b32d499e8b9482c04030f5f5_master.jpg','2025-09-28 18:34:00',''),(59,'https://www.bluestone.com.vn/blogs/vao-bep/chien-khoai-tay-bang-noi-chien-khong-dau?srsltid=AfmBOoqJ8H9yAh7LTcp--k4VJ6EP2X2OuqePxFBm3Y8K1ygpG_FpQDcN','2025-09-28 18:34:00',''),(60,'https://img.giftpop.vn/brand/LOTTERIA/1PEMP2010312501_BASIC_origin.jpg','2025-09-28 18:34:00',''),(61,'https://cdn.tgdd.vn/Files/2020/03/02/1239549/2-cong-thuc-lam-banh-hotdog-xuc-xich-hotdog-pho-mai-han-quoc-gay-nghien-14-760x367.png','2025-09-28 18:34:00',''),(62,'https://checkinvietnam.vtc.vn/media/20211221/files/pizza-xuc-xich-pho-mai-vuong.jpg','2025-09-28 18:34:00',''),(63,'https://i.ytimg.com/vi/ng3vo1RmeyQ/maxresdefault.jpg','2025-09-28 18:34:00',''),(64,'https://storage.googleapis.com/onelife-public/blog.onelife.vn/2021/10/cach-lam-banh-mi-sandwich-trung-jambon-mon-an-sang-349515833958.jpg','2025-09-28 18:34:00',''),(65,'https://cdnv2.tgdd.vn/bhx-static/bhx/Products/Images/7259/332717/bhx/frame-3475095-2-1_202412022211137044.jpg','2025-09-28 18:34:00',''),(66,'https://cdn.tgdd.vn/Files/2022/03/07/1418886/9-cach-lam-salad-tron-mayonnaise-giam-can-tai-nha-hieu-qua-202203071357195806.jpg','2025-09-28 18:34:00',''),(68,'https://res.cloudinary.com/dytix5ybu/image/upload/v1759138848/product/TraDaoCamSa.jpg.jpg','2025-09-29 17:03:26','41'),(69,'https://res.cloudinary.com/dytix5ybu/image/upload/v1759138848/product/TraDaoCamSa.jpg.jpg','2025-10-09 03:37:09','42'),(70,'https://res.cloudinary.com/dytix5ybu/image/upload/v1759138848/product/TraDaoCamSa.jpg.jpg','2025-10-09 03:39:58','43'),(71,'https://res.cloudinary.com/dytix5ybu/image/upload/v1763459988/product/thum-1200x676-16.jpg.jpg','2025-11-18 09:59:49','product/thum-1200x676-16.jpg'),(72,'https://res.cloudinary.com/dytix5ybu/image/upload/v1766413249/product/t1.jpg.jpg','2025-12-22 14:20:49','product/t1.jpg');
/*!40000 ALTER TABLE `Images` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ProductImages`
--

DROP TABLE IF EXISTS `ProductImages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ProductImages` (
  `product_id` int NOT NULL,
  `image_id` int NOT NULL,
  `is_main` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`product_id`,`image_id`),
  KEY `fk_image` (`image_id`),
  CONSTRAINT `fk_image` FOREIGN KEY (`image_id`) REFERENCES `Images` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_product` FOREIGN KEY (`product_id`) REFERENCES `Products` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ProductImages`
--

LOCK TABLES `ProductImages` WRITE;
/*!40000 ALTER TABLE `ProductImages` DISABLE KEYS */;
INSERT INTO `ProductImages` VALUES (16,37,1),(16,38,0),(16,39,0),(19,46,1),(19,47,0),(20,48,1),(20,49,0),(20,50,0),(21,51,1),(21,52,0),(22,53,1),(23,54,1),(24,55,1),(25,56,1),(26,57,1),(27,58,1),(29,60,1),(30,61,1),(31,62,1),(32,63,1),(33,64,1),(34,65,1),(35,66,1),(41,68,1),(42,69,1),(43,70,1),(44,71,1),(45,72,1);
/*!40000 ALTER TABLE `ProductImages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Products`
--

DROP TABLE IF EXISTS `Products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Products` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text,
  `price` decimal(10,2) NOT NULL,
  `qty_initial` int DEFAULT '0',
  `qty_sold` int DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  CONSTRAINT `Products_chk_1` CHECK ((`price` >= 0)),
  CONSTRAINT `Products_chk_2` CHECK ((`qty_initial` >= 0)),
  CONSTRAINT `Products_chk_3` CHECK ((`qty_sold` >= 0))
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Products`
--

LOCK TABLES `Products` WRITE;
/*!40000 ALTER TABLE `Products` DISABLE KEYS */;
INSERT INTO `Products` VALUES (16,'Nước Cam Ép','Nước cam tươi nguyên chất',15000.00,100,35,'2025-09-24 02:46:55','2025-09-24 09:46:55'),(19,'Sữa Tươi','Sữa tươi tiệt trùng nguyên chất',18000.00,80,30,'2025-09-24 02:46:55','2025-09-24 09:46:55'),(20,'Sữa Đậu Nành','Thức uống từ đậu nành bổ dưỡng',12000.00,90,22,'2025-09-24 02:46:55','2025-09-24 09:46:55'),(21,'Trà Sữa Trân Châu','Trà sữa kèm trân châu dai ngon',35000.00,100,50,'2025-09-24 02:46:55','2025-09-24 09:46:55'),(22,'Trà Đào Cam Sả','Trà đào cam sả mát lạnh',30000.00,60,20,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(23,'Cà Phê Đen','Cà phê đen nguyên chất',20000.00,80,35,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(24,'Cà Phê Sữa','Cà phê sữa đá truyền thống',25000.00,90,40,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(25,'Sinh Tố Bơ','Sinh tố bơ béo ngậy',40000.00,50,18,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(26,'Burger Bò Phô Mai','Bánh burger bò kèm phô mai tan chảy',45000.00,50,20,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(27,'Burger Gà Giòn','Bánh burger gà chiên giòn',40000.00,60,25,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(29,'Gà Rán 2 Miếng','Gà rán giòn rụm, hương vị đặc trưng',60000.00,80,30,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(30,'Hotdog Xúc Xích','Bánh mì kẹp xúc xích và tương cà',30000.00,70,20,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(31,'Pizza Phô Mai','Pizza nhỏ phủ phô mai mozzarella',70000.00,40,15,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(32,'Pizza Hải Sản','Pizza hải sản tươi ngon',85000.00,35,10,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(33,'Sandwich Thịt Nguội','Bánh sandwich kẹp thịt nguội và rau',35000.00,60,22,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(34,'Mì Ý Sốt Bò Bằm','Mì Ý sốt cà chua bò bằm',65000.00,45,18,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(35,'Salad Rau Trộn','Salad rau củ tươi mát',30000.00,50,12,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(41,'Trà Đào Cam Sả','Trà đào cam xả 100% làm từ thiên nhiên',25000.00,100,0,'2025-09-29 10:03:24','2025-09-29 17:03:24'),(42,'Trà đào cam sả','Ngon bổ rẻ khỏe',20000.00,100,5,'2025-10-08 20:37:08','2025-10-09 03:37:07'),(43,'Trà Đào Cam Sả 04','ngon ngon ngon',12300.00,23,0,'2025-10-08 20:39:57','2025-10-09 03:39:57'),(44,'Trà chanh nhiệt đới','Ngon nhức nách',30000.00,123,0,'2025-11-18 02:59:50','2025-11-18 09:59:49'),(45,'Tra chanh sieu ngon','bla bla bla',23000.00,100,0,'2025-12-22 07:20:50','2025-12-22 14:20:49');
/*!40000 ALTER TABLE `Products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ReviewImages`
--

DROP TABLE IF EXISTS `ReviewImages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ReviewImages` (
  `review_id` int NOT NULL,
  `image_id` int NOT NULL,
  UNIQUE KEY `uq_review_image` (`review_id`,`image_id`),
  KEY `fk_reviewimages_image` (`image_id`),
  CONSTRAINT `fk_reviewimages_image` FOREIGN KEY (`image_id`) REFERENCES `Images` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_reviewimages_review` FOREIGN KEY (`review_id`) REFERENCES `Reviews` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ReviewImages`
--

LOCK TABLES `ReviewImages` WRITE;
/*!40000 ALTER TABLE `ReviewImages` DISABLE KEYS */;
/*!40000 ALTER TABLE `ReviewImages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Reviews`
--

DROP TABLE IF EXISTS `Reviews`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Reviews` (
  `id` int NOT NULL AUTO_INCREMENT,
  `product_id` int NOT NULL,
  `user_id` int NOT NULL,
  `order_id` int NOT NULL,
  `rate` int NOT NULL,
  `content` varchar(1000) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_user_product_order` (`user_id`,`product_id`,`order_id`),
  KEY `fk_review_product` (`product_id`),
  KEY `fk_review_order` (`order_id`),
  CONSTRAINT `fk_review_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
  CONSTRAINT `fk_review_product` FOREIGN KEY (`product_id`) REFERENCES `Products` (`id`),
  CONSTRAINT `fk_review_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `Reviews_chk_1` CHECK ((`rate` between 1 and 5))
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Reviews`
--

LOCK TABLES `Reviews` WRITE;
/*!40000 ALTER TABLE `Reviews` DISABLE KEYS */;
INSERT INTO `Reviews` VALUES (4,16,37,7,4,'ngon vo cung luon','2025-10-08 20:04:33');
/*!40000 ALTER TABLE `Reviews` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `messages`
--

DROP TABLE IF EXISTS `messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `messages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_id` int NOT NULL,
  `sender_id` int NOT NULL,
  `receiver_id` int NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `is_read` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  KEY `fk_order` (`order_id`),
  KEY `fk_sender` (`sender_id`),
  KEY `fk_receiver` (`receiver_id`),
  CONSTRAINT `fk_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
  CONSTRAINT `fk_receiver` FOREIGN KEY (`receiver_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_sender` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=61 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messages`
--

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;
INSERT INTO `messages` VALUES (5,6,39,37,'Đơn hàng đã sẵn sàng giao','2025-09-30 11:15:25',0),(6,6,39,37,'Đơn hàng đã sẵn sàng giao','2025-09-30 16:33:21',0),(7,6,39,37,'chao lan 1','2025-09-30 16:33:38',0),(8,6,39,37,'chao lan 2','2025-09-30 16:33:42',0),(9,6,39,37,'chao lan 2','2025-09-30 16:33:45',0),(10,6,37,39,'chao lai lan 1','2025-09-30 16:34:07',0),(11,6,39,37,'hello customer','2025-09-30 16:46:58',0),(12,6,37,39,'hi shipper','2025-09-30 16:47:20',0),(13,6,39,37,'chao lan 2','2025-10-03 14:52:12',0),(14,6,37,39,'chao shipper','2025-10-03 14:53:45',0),(16,6,39,37,'hi','2025-10-09 02:27:28',0),(17,6,39,37,'hi','2025-10-09 02:27:39',0),(18,6,39,37,'hi','2025-10-09 02:42:40',0),(19,7,37,39,'chao lan 2','2025-10-09 02:42:48',0),(20,6,39,37,'hi','2025-10-09 02:44:22',0),(21,6,39,37,'hi','2025-10-09 02:44:58',0),(22,6,39,37,'hi','2025-10-09 02:45:36',0),(23,7,37,39,'chao lan 2','2025-10-09 02:47:05',0),(24,7,37,39,'chao lan 2','2025-10-09 02:47:24',0),(25,6,39,37,'hi','2025-10-09 03:44:55',0),(26,6,39,37,'hi','2025-10-09 03:48:31',0),(27,6,39,37,'hi','2025-10-09 03:48:58',0),(28,6,39,37,'hi','2025-10-09 04:02:08',0),(29,6,39,37,'hi','2025-10-09 04:03:24',0),(30,7,37,39,'chao lan 2','2025-10-09 04:03:34',0),(31,6,39,37,'hi','2025-10-09 04:03:39',0),(32,14,39,37,'alo','2025-10-09 04:46:49',0),(33,14,39,37,'hu\ndfdsa','2025-10-09 04:47:46',0),(34,6,37,39,'chao lan 2','2025-10-09 04:50:48',0),(35,6,37,39,'chao lan 2','2025-10-09 04:50:48',0),(36,6,37,39,'chao lan 2','2025-10-09 04:50:49',0),(37,6,37,39,'chao lan 2','2025-10-09 04:50:50',0),(38,6,37,39,'chao lan 2','2025-10-09 04:50:51',0),(39,6,37,39,'chao lan 2','2025-10-09 04:50:51',0),(40,6,37,39,'chao lan 2','2025-10-09 04:50:51',0),(41,6,37,39,'chao lan 2','2025-10-09 04:52:05',0),(42,6,37,39,'chao lan 2','2025-10-09 04:52:06',0),(43,6,39,37,'chao lan 2','2025-10-09 04:55:56',0),(44,6,37,39,'chao lan 2','2025-10-09 04:56:08',0),(45,6,39,37,'chao lan 43453','2025-10-09 04:57:04',0),(46,6,39,37,'chao lan 43453','2025-10-09 04:59:43',0),(47,6,39,37,'chao lan 43453dsf','2025-10-09 05:00:40',0),(48,14,39,37,'chaof cau','2025-10-09 05:01:03',0),(49,14,39,37,'hello','2025-10-09 05:01:21',0),(50,6,37,39,'hello shipper','2025-10-09 05:01:53',0),(51,6,37,39,'hello shipper, howw are you','2025-10-09 05:25:33',0),(52,6,37,39,'hello shipper, howw are you','2025-10-09 05:26:01',0),(53,14,39,37,'as','2025-10-09 05:26:12',0),(54,14,39,37,'sdsadsd','2025-10-09 05:26:18',0),(55,14,39,37,'sd','2025-10-09 05:26:51',0),(56,6,37,39,'hello shipper, howw are you','2025-10-09 05:26:56',0),(57,14,39,37,'sds','2025-10-09 05:27:00',0),(58,14,39,37,'sds','2025-10-09 05:44:27',0),(59,6,37,39,'hello shipper, howw are you','2025-10-09 05:44:35',0),(60,14,39,37,'chao cau cau co khoe khong cut ra nhan hang gium tao\n','2025-10-09 05:45:58',0);
/*!40000 ALTER TABLE `messages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `order_items`
--

DROP TABLE IF EXISTS `order_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `order_items` (
  `id` int NOT NULL AUTO_INCREMENT,
  `order_id` int NOT NULL,
  `product_id` int NOT NULL,
  `quantity` int NOT NULL,
  `price` decimal(10,2) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `order_id` (`order_id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `order_items_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE,
  CONSTRAINT `order_items_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `Products` (`id`) ON DELETE CASCADE,
  CONSTRAINT `order_items_chk_1` CHECK ((`quantity` > 0)),
  CONSTRAINT `order_items_chk_2` CHECK ((`price` >= 0))
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `order_items`
--

LOCK TABLES `order_items` WRITE;
/*!40000 ALTER TABLE `order_items` DISABLE KEYS */;
INSERT INTO `order_items` VALUES (8,7,16,7,15000.00),(9,8,34,6,65000.00),(10,8,35,7,30000.00),(13,10,29,1,60000.00),(14,10,27,2,40000.00),(15,11,26,2,45000.00),(16,11,30,2,30000.00),(17,12,26,2,45000.00),(18,12,30,2,30000.00),(19,13,26,2,45000.00),(20,13,30,2,30000.00),(21,14,26,2,45000.00),(22,14,30,2,30000.00);
/*!40000 ALTER TABLE `order_items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `shipper_id` int DEFAULT NULL,
  `payment_status` enum('unpaid','paid','refunded') DEFAULT 'unpaid',
  `order_status` enum('pending','processing','shipping','delivered','cancelled') DEFAULT 'pending',
  `latitude` decimal(10,8) NOT NULL,
  `longitude` decimal(11,8) NOT NULL,
  `total_amount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `thumbnail_id` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_orders_thumbnail` (`thumbnail_id`),
  CONSTRAINT `fk_orders_thumbnail` FOREIGN KEY (`thumbnail_id`) REFERENCES `Images` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orders`
--

LOCK TABLES `orders` WRITE;
/*!40000 ALTER TABLE `orders` DISABLE KEYS */;
INSERT INTO `orders` VALUES (6,37,39,'unpaid','processing',21.02851100,105.80481700,149000.00,28,'2025-09-27 02:01:50','2025-12-22 13:57:05'),(7,37,39,'unpaid','processing',21.02851100,105.80481700,165000.00,34,'2025-09-27 03:21:13','2025-12-22 13:57:05'),(8,37,NULL,'unpaid','processing',21.02851100,105.80481700,600000.00,65,'2025-10-02 03:31:43','2025-12-22 13:57:05'),(10,37,NULL,'unpaid','processing',21.02851100,105.80481700,140000.00,60,'2025-10-02 03:33:05','2025-12-22 13:57:05'),(11,37,NULL,'unpaid','processing',21.02851100,105.80481700,150000.00,57,'2025-10-02 03:33:18','2025-12-22 13:57:05'),(12,37,NULL,'unpaid','processing',21.02851100,105.80481700,150000.00,57,'2025-10-02 03:33:19','2025-12-22 13:57:05'),(13,37,NULL,'unpaid','processing',21.02851100,105.80481700,150000.00,57,'2025-10-02 03:33:20','2025-12-22 13:57:05'),(14,37,39,'unpaid','processing',21.02851100,105.80481700,150000.00,57,'2025-10-08 19:50:51','2025-12-22 13:57:05');
/*!40000 ALTER TABLE `orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `refresh_tokens`
--

DROP TABLE IF EXISTS `refresh_tokens`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `refresh_tokens` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `token` text NOT NULL,
  `expires_at` timestamp NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `refresh_tokens_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=106 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `refresh_tokens`
--

LOCK TABLES `refresh_tokens` WRITE;
/*!40000 ALTER TABLE `refresh_tokens` DISABLE KEYS */;
INSERT INTO `refresh_tokens` VALUES (94,36,'xHbaeye3kmd4a05Ojlqovl6Da-gTuUVaA_MBdfiNmxI=','2025-11-19 08:36:13','2025-11-12 15:36:12','2025-11-12 15:36:12'),(95,36,'Cv7q48Pb9sC4UW-VpjlC_Uw4MgwQn_NGEsh8Yd0Ui5Y=','2025-11-24 23:50:26','2025-11-18 06:50:26','2025-11-18 06:50:26'),(96,36,'L-q0B0gGxf4-sGbVyBL8lgdRYOiNEaV4IsGZ5E2fkyE=','2025-11-25 00:02:14','2025-11-18 07:02:14','2025-11-18 09:59:46'),(97,36,'bcgnlkKUoaQ9QohBHYiby1iF3mlr-9vZ3is7qxSM9hc=','2025-11-25 09:58:13','2025-11-18 16:58:12','2025-11-18 16:58:12'),(102,36,'xkk-e-7Vlm-ap5Qhpt4-7NwO5rx82qZgjYKJn37kuXc=','2025-12-28 22:33:12','2025-12-22 05:33:12','2025-12-22 05:33:12'),(103,36,'JR5nfMdOAcpppCx1X0MEY8m3ortUBIcde_OfjcJcPx8=','2025-12-29 02:36:52','2025-12-22 09:36:52','2025-12-22 09:36:52'),(104,39,'WPsCL0g_rmHvU4VjfirXIR2KXWbFshCsnPLfnlrupkU=','2025-12-29 06:55:28','2025-12-22 13:55:27','2025-12-22 13:55:27'),(105,36,'J8ABBeCKYXXQTtWdDyhbl-RZL3pD5NsM0ggRkZ7JFjk=','2025-12-29 07:20:10','2025-12-22 14:20:09','2025-12-22 14:20:09');
/*!40000 ALTER TABLE `refresh_tokens` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` text NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `role` enum('customer','shipper','supplier','admin') NOT NULL DEFAULT 'customer',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `otp_code` varchar(10) DEFAULT NULL,
  `reset_otp_num` int DEFAULT '0',
  `otp_expires_at` timestamp NULL DEFAULT NULL,
  `reset_otp` varchar(10) DEFAULT NULL,
  `reset_otp_expires_at` timestamp NULL DEFAULT NULL,
  `status` tinyint DEFAULT '0' COMMENT '0=inactive, 1=active, 2=banned, 3=suspended',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (35,'Admin','admin@example.com','$2a$10$OeGOOr1OBOm6VmxnCJdbOej0325iqiiSQpMtQdQaGSXF7DbRRcmQe','0000000000','Admin Address','admin','2025-09-24 09:36:35',NULL,0,NULL,NULL,NULL,1),(36,'Admin','admin@gmail.com','$2a$10$KG.aMLHJSAHMb2QwneUVwOZ2zISBuLm5K.15hTT9O6DC.o6CQmbHy','','','admin','2025-09-24 02:37:46',NULL,0,NULL,NULL,NULL,1),(37,'User1','user1@gmail.com','$2a$10$s50wUtykgFKNELyjFf2z9.PIO82CteRVDu9ATLmOgj0ns0doMOMqW','0987777777','Ha Noi','customer','2025-09-26 00:40:40','529084',0,'2025-09-26 00:50:40',NULL,NULL,1),(38,'admin01','admin1@gmail.com','$2a$10$o8m4W5wobYd10hBL/uemL.F0OtqY0JKGZk6.0cb3dkiy7M4tpOlNu','0987777777','Ha Noi','admin','2025-09-27 01:53:09',NULL,0,NULL,NULL,NULL,1),(39,'Shipper1','shipper1@gmail.com','$2a$10$flof0WS10vwudJj7394emO4v4ZgIp9HAxi3Y6Gs0edU4mP11tPs3i','0987876768','Ha Noi','shipper','2025-09-27 10:23:26','563546',0,'2025-09-28 04:21:28',NULL,NULL,1),(42,'Shipper2','shipper2@gmail.com','$2a$10$IImb9LsKaz02AA4L1jPEYuYpaUtYT5G7Ly/V6F2WkfMcdk6KdwG5m','0987654321','Ha Dong','shipper','2025-10-03 01:29:39',NULL,0,NULL,NULL,NULL,1),(43,'Shipper3','Shipper3@gmail.com','$2a$10$ccyZzkfkYm/wdZ09p3kT4O.lcgHZiRAJNZsI1VKSO6xZX6duyZhFS','0997678942','Hoa Binh','shipper','2025-10-03 01:36:41',NULL,0,NULL,NULL,NULL,1),(44,'Shipper4','shipper4@gmail.com','$2a$10$0Ujgcug5LjUmDamt1TlcMOx6bxvDec1PYm8uTWFG8lTORESZNRjYq','0988897655','Ha Noi','shipper','2025-10-03 07:37:24',NULL,0,NULL,NULL,NULL,1),(46,'Nguyen Duc Anh 02','nguyenduca03@gmail.com','$2a$10$OGq.I4HfLqsdDCh7hvFUW.F.MJam3oBtDcvpVeQnO49/pAYZzy6bO','0976673117','Ha Noi','customer','2025-10-03 19:11:30',NULL,0,NULL,NULL,NULL,1),(47,'Duck Anh','darke7824@gmail.com','$2a$10$doCT9YyULFXM.byN8ISN6eG2itWjcOPdicRXN1FTE0PaLVjefI/6e','0987777772','Ha Noi','customer','2025-11-04 09:13:12','911611',0,'2025-11-04 09:23:12',NULL,NULL,1),(48,'Shipper5','shipper5@gmail.com','$2a$10$Eg3IJD7ch3v285sEytbK0.DDdVPt4Vr4rvOho5jhUe1SyMd8.6d5K','0988888888','Hung Dao - Ha Noi','shipper','2025-11-17 23:59:31',NULL,0,NULL,NULL,NULL,1),(49,'Shipper6','shipper6@gmail.com','$2a$10$I7QZ0ADUXxdwSkseMvAu2epE5j8UbzRTiY2bQxP8gklEp3f2Kuok.','0987654321','Ha Tay','shipper','2025-12-22 02:37:40',NULL,0,NULL,NULL,NULL,1);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-12-22 21:24:18
