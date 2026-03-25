-- phpMyAdmin SQL Dump
-- version 5.2.2deb1+deb13u1
-- https://www.phpmyadmin.net/
--
-- Servidor: localhost:3306
-- Tiempo de generación: 25-03-2026 a las 17:17:06
-- Versión del servidor: 11.8.3-MariaDB-0+deb13u1 from Debian
-- Versión de PHP: 8.4.16

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Base de datos: `Tocke`
--

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `categorias`
--

CREATE TABLE `categorias` (
  `id_cat` int(4) NOT NULL,
  `nombre` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `categorias`
--

INSERT INTO `categorias` (`id_cat`, `nombre`) VALUES
(2, 'Sanguches'),
(5, 'Papas'),
(6, 'Postres'),
(7, 'Pizzas'),
(8, 'Chorrillanas'),
(10, 'Helados'),
(11, 'empanadas'),
(12, 'Hamburguesas'),
(13, 'Promociones'),
(15, 'Bebidas');

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `configuracion`
--

CREATE TABLE `configuracion` (
  `clave` varchar(50) NOT NULL,
  `valor` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `configuracion`
--

INSERT INTO `configuracion` (`clave`, `valor`) VALUES
('pedidos_online', 'cerrado');

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `ingredientes`
--

CREATE TABLE `ingredientes` (
  `id_ing` int(11) NOT NULL,
  `nombre` varchar(50) NOT NULL,
  `stock` int(11) NOT NULL DEFAULT 0,
  `stock_minimo` int(11) NOT NULL DEFAULT 20,
  `unidad` varchar(20) NOT NULL DEFAULT 'unidad'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `ingredientes`
--

INSERT INTO `ingredientes` (`id_ing`, `nombre`, `stock`, `stock_minimo`, `unidad`) VALUES
(7, 'Pan de Hamburguesa', 100, 90, 'unidad'),
(8, 'Pan de Hot dog', 50, 56, 'unidad'),
(10, 'Envase plumavit', 97, 10, 'unidad'),
(11, 'Envase carton', 100, 10, 'unidad'),
(14, 'Bebidas', 11, 10, 'unidad');

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `mesas`
--

CREATE TABLE `mesas` (
  `id_mesa` int(11) NOT NULL,
  `nombre` varchar(20) NOT NULL,
  `ocupada` tinyint(1) DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `mesas`
--

INSERT INTO `mesas` (`id_mesa`, `nombre`, `ocupada`) VALUES
(1, 'Mesa 1', 0),
(2, 'Mesa 2', 0),
(3, 'Mesa 3', 0),
(4, 'Barra 1', 0),
(5, 'Barra 2', 0),
(6, 'Barra 3', 0),
(7, 'Barra 4', 0),
(8, 'Barra 5', 0),
(9, 'Barra 6', 0),
(10, 'Barra 7', 0),
(11, 'Barra 8', 0);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `modificadores`
--

CREATE TABLE `modificadores` (
  `id_mod` int(11) NOT NULL,
  `id_pro` int(11) NOT NULL,
  `nombre` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `modificadores`
--

INSERT INTO `modificadores` (`id_mod`, `id_pro`, `nombre`) VALUES
(1, 8, 'sin mayo'),
(2, 8, 'sin tomate'),
(3, 2, 'sin mayonesa'),
(6, 22, 'sin mayo'),
(7, 23, 'ChurrascoItaliano'),
(8, 23, 'Coca cola'),
(9, 23, 'Fanta'),
(11, 23, 'empanadas de queso');

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `pedidos`
--

CREATE TABLE `pedidos` (
  `id_ped` int(4) NOT NULL,
  `fecha` datetime NOT NULL,
  `total` int(4) NOT NULL,
  `cliente` varchar(50) DEFAULT NULL,
  `tipo_pedido` varchar(20) NOT NULL DEFAULT 'servirse',
  `id_mesa` int(11) DEFAULT NULL,
  `estado` varchar(20) NOT NULL DEFAULT 'abierto'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `pedidos`
--

INSERT INTO `pedidos` (`id_ped`, `fecha`, `total`, `cliente`, `tipo_pedido`, `id_mesa`, `estado`) VALUES
(25, '2026-03-17 07:44:01', 800, NULL, 'Servir', NULL, 'cerrado'),
(26, '2026-03-17 07:45:00', 13100, NULL, 'Servir', NULL, 'cerrado'),
(27, '2026-03-17 07:54:46', 22300, NULL, 'Servir', NULL, 'cerrado'),
(28, '2026-03-17 07:55:17', 6600, NULL, 'Servir', NULL, 'cerrado'),
(29, '2026-03-17 07:56:01', 4600, NULL, 'Servir', NULL, 'cerrado'),
(30, '2026-03-17 07:57:59', 16600, NULL, 'Servir', NULL, 'cerrado'),
(31, '2026-03-17 08:01:55', 15700, NULL, 'Servir', NULL, 'cerrado'),
(32, '2026-03-17 08:14:33', 12000, NULL, 'Servir', NULL, 'cerrado'),
(33, '2026-03-17 12:32:36', 26100, NULL, 'Servir', NULL, 'cerrado'),
(34, '2026-03-17 12:33:29', 15700, NULL, 'Servir', NULL, 'cerrado'),
(35, '2026-03-17 12:34:01', 35700, NULL, 'Servir', NULL, 'cerrado'),
(36, '2026-03-17 12:46:48', 12800, NULL, 'Servir', NULL, 'cerrado'),
(37, '2026-03-17 12:49:41', 11700, NULL, 'Servir', NULL, 'cerrado'),
(38, '2026-03-17 13:11:59', 14800, NULL, 'Servir', NULL, 'cerrado'),
(39, '2026-03-17 13:17:14', 13000, NULL, 'Servir', NULL, 'cerrado'),
(40, '2026-03-17 13:28:32', 3500, NULL, 'Servir', NULL, 'cerrado'),
(41, '2026-03-17 13:29:06', 3500, NULL, 'Servir', NULL, 'cerrado'),
(42, '2026-03-17 13:33:59', 3500, NULL, 'Servir', NULL, 'cerrado'),
(43, '2026-03-17 13:35:40', 7000, NULL, 'Servir', NULL, 'cerrado'),
(44, '2026-03-17 13:39:48', 18500, NULL, 'Servir', NULL, 'cerrado'),
(45, '2026-03-17 13:40:35', 18000, NULL, 'Servir', NULL, 'cerrado'),
(46, '2026-03-17 13:57:45', 19800, NULL, 'Servir', NULL, 'cerrado'),
(47, '2026-03-17 14:09:46', 13800, NULL, 'Servir', NULL, 'cerrado'),
(48, '2026-03-17 15:54:46', 17300, NULL, 'Servir', NULL, 'cerrado'),
(49, '2026-03-17 16:12:23', 7300, NULL, 'Servir', NULL, 'cerrado'),
(50, '2026-03-18 08:28:13', 3500, 'Juanito', 'Servir', NULL, 'cerrado'),
(51, '2026-03-18 08:34:10', 1700, 'Paula', 'Servir', NULL, 'cerrado'),
(52, '2026-03-18 17:58:51', 900, 'dsfsf', 'Servir', NULL, 'cerrado'),
(53, '2026-03-18 17:59:09', 900, 'dsfsf', 'Servir', NULL, 'cerrado'),
(54, '2026-03-20 16:50:49', 10500, 'sdf', 'Servir', NULL, 'cerrado'),
(55, '2026-03-20 17:26:47', 4400, 'fdgd', 'Servir', NULL, 'cerrado'),
(56, '2026-03-21 09:52:14', 900, 'adfaf', 'retiro', NULL, 'cerrado'),
(57, '2026-03-21 09:52:24', 900, 'adfaf', 'retiro', NULL, 'cerrado'),
(58, '2026-03-21 09:59:41', 1800, 'cxcx', 'retiro', NULL, 'cerrado'),
(59, '2026-03-21 10:02:21', 900, 'erer', 'Llevar', NULL, 'cerrado'),
(60, '2026-03-21 10:13:15', 1500, 'dfsdf', 'Retiro', NULL, 'cerrado'),
(61, '2026-03-21 10:14:06', 12000, 'kjl', 'Delivery', NULL, 'cerrado'),
(62, '2026-03-21 10:27:37', 7500, 'hh', 'Retiro', NULL, 'cerrado'),
(63, '2026-03-21 10:28:36', 4500, 'b', 'Servir', NULL, 'cerrado'),
(64, '2026-03-21 13:30:46', 1500, 'as', 'Llevar', NULL, 'cerrado'),
(65, '2026-03-21 13:33:52', 3000, 'sd', 'Servir', NULL, 'cerrado'),
(66, '2026-03-21 13:34:34', 3000, 'sd', 'Servir', NULL, 'cerrado'),
(67, '2026-03-21 13:37:16', 1500, 'sda', 'Retiro', NULL, 'cerrado'),
(68, '2026-03-21 13:43:54', 1500, 'nb', 'Retiro', NULL, 'cerrado'),
(69, '2026-03-21 13:44:16', 3600, 'nn', 'Delivery', NULL, 'cerrado'),
(70, '2026-03-21 13:48:48', 3600, 'dd', 'Retiro', NULL, 'cerrado'),
(71, '2026-03-21 13:59:08', 30000, 'dtgh', 'Servir', NULL, 'cerrado'),
(72, '2026-03-21 14:25:14', 13300, 'jhg', 'Retiro', NULL, 'cerrado'),
(73, '2026-03-21 14:33:43', 3600, 'uu', 'Servir', NULL, 'cerrado'),
(74, '2026-03-21 14:34:28', 15000, 'ewrwer', 'Delivery', NULL, 'cerrado'),
(75, '2026-03-21 14:38:41', 14500, 'ssss', 'Llevar', NULL, 'cerrado'),
(76, '2026-03-21 14:40:34', 5400, 'erer', 'Retiro', NULL, 'cerrado'),
(77, '2026-03-21 14:43:32', 8400, 'cc', 'Llevar', NULL, 'cerrado'),
(78, '2026-03-21 14:45:24', 2000, 'adfaf', 'Retiro', NULL, 'cerrado'),
(79, '2026-03-21 14:46:31', 4500, 'cxcx', 'Llevar', NULL, 'cerrado'),
(80, '2026-03-21 14:55:18', 5400, 'kjl', 'Llevar', NULL, 'cerrado'),
(81, '2026-03-21 14:57:51', 6000, 'kjl', 'Llevar', NULL, 'cerrado'),
(82, '2026-03-21 14:58:52', 3600, 'dfsdf', 'Retiro', NULL, 'cerrado'),
(83, '2026-03-21 15:03:00', 6500, 'erer', 'Delivery', NULL, 'cerrado'),
(84, '2026-03-21 15:06:03', 4000, 'cvcvc', 'Retiro', NULL, 'cerrado'),
(85, '2026-03-21 15:12:04', 3000, 'dfsdf', 'Retiro', NULL, 'cerrado'),
(86, '2026-03-21 15:12:29', 4000, 'kjl', 'Delivery', NULL, 'cerrado'),
(87, '2026-03-21 15:17:59', 9900, 'erer', 'Retiro', NULL, 'cerrado'),
(88, '2026-03-21 16:24:55', 12900, 'kjl', 'Llevar', NULL, 'cerrado'),
(89, '2026-03-21 16:25:19', 6800, 'hjhjj', 'Llevar', NULL, 'cerrado'),
(90, '2026-03-21 16:26:49', 9000, 'dfsdf', 'Llevar', NULL, 'cerrado'),
(91, '2026-03-21 16:33:27', 16000, 'cxcx', 'Servir', NULL, 'cerrado'),
(92, '2026-03-21 16:50:29', 2700, 'adfaf', 'Retiro', NULL, 'cerrado'),
(93, '2026-03-21 17:15:38', 4200, 'adfaf', 'Retiro', NULL, 'cerrado'),
(94, '2026-03-21 17:39:26', 3300, 'dfsdf', 'Servir', 2, 'cerrado'),
(95, '2026-03-21 17:40:23', 2000, 'zczx', 'Servir', 7, 'cerrado'),
(96, '2026-03-21 17:40:55', 3500, 'cxcx', 'Servir', 10, 'cerrado'),
(97, '2026-03-21 17:53:08', 1500, 'cxcx', 'Servir', 11, 'cerrado'),
(98, '2026-03-21 18:02:56', 2000, 'adfaf', 'Servir', 9, 'cerrado'),
(99, '2026-03-22 18:42:10', 10000, 'erer', 'Retiro', NULL, 'cerrado'),
(100, '2026-03-24 10:28:43', 3500, 'erer', 'Ir comiendo', NULL, 'cerrado'),
(101, '2026-03-24 10:40:07', 7000, 'adfaf', 'Llevar', NULL, 'cerrado'),
(102, '2026-03-24 14:20:42', 4000, 'cxcx', 'Delivery', NULL, 'cerrado'),
(103, '2026-03-24 14:23:43', 6000, 'cxcx', 'Servir', NULL, 'cerrado'),
(104, '2026-03-24 19:05:41', 2000, 'kjl', 'Servir', NULL, 'cerrado'),
(105, '2026-03-24 19:10:45', 1800, 'erer', 'Servir', 10, 'cerrado'),
(106, '2026-03-25 09:13:24', 6000, 'erer', 'Retiro', NULL, 'cerrado'),
(107, '2026-03-25 09:30:01', 5000, 'adfaf', 'Ir comiendo', NULL, 'cerrado'),
(108, '2026-03-25 09:33:11', 3000, 'dfsdf', 'Delivery', NULL, 'cerrado'),
(109, '2026-03-25 09:34:49', 7000, 'erer', 'Llevar', NULL, 'cerrado'),
(110, '2026-03-25 10:00:17', 10000, 'kjl', 'Ir comiendo', NULL, 'cerrado'),
(111, '2026-03-25 10:11:03', 5000, 'sdf', 'Retiro', NULL, 'cerrado'),
(112, '2026-03-25 10:11:34', 6000, 'dfsdf', 'Servir', 8, 'cerrado'),
(113, '2026-03-25 10:16:31', 2000, 'sdf', 'Retiro', NULL, 'cerrado'),
(114, '2026-03-25 10:34:26', 14000, 'sgf', 'Llevar', NULL, 'cerrado'),
(115, '2026-03-25 10:36:54', 3500, 'adfaf', 'Ir comiendo', NULL, 'cerrado'),
(116, '2026-03-25 10:43:50', 10000, 'kjl', 'Llevar', NULL, 'cerrado'),
(117, '2026-03-25 11:08:11', 3500, 'cc', 'Delivery', NULL, 'cerrado'),
(118, '2026-03-25 12:39:05', 11000, 'Paula', 'Retiro', NULL, 'cerrado'),
(119, '2026-03-25 14:50:11', 6000, 'dfsdf', 'Servir', 7, 'abierto'),
(120, '2026-03-25 15:04:39', 5000, 'cxcx', 'Llevar', NULL, 'abierto');

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `pedidos_detalle`
--

CREATE TABLE `pedidos_detalle` (
  `id_ped` int(4) NOT NULL,
  `id_pro` int(4) NOT NULL,
  `cantidad` int(4) NOT NULL DEFAULT 1,
  `precio` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `pedidos_detalle`
--

INSERT INTO `pedidos_detalle` (`id_ped`, `id_pro`, `cantidad`, `precio`) VALUES
(26, 8, 1, 3500),
(26, 2, 1, 3000),
(27, 8, 1, 3500),
(27, 2, 1, 3000),
(27, 7, 1, 2800),
(30, 2, 3, 3000),
(31, 8, 3, 3500),
(33, 8, 3, 3500),
(33, 2, 2, 3000),
(34, 8, 3, 3500),
(35, 8, 2, 3500),
(35, 2, 2, 3000),
(35, 7, 2, 2800),
(36, 12, 1, 2000),
(37, 8, 1, 3500),
(37, 2, 1, 3000),
(38, 8, 2, 3500),
(39, 8, 2, 3500),
(40, 8, 1, 3500),
(41, 8, 1, 3500),
(42, 8, 1, 3500),
(43, 8, 2, 3500),
(44, 8, 1, 3500),
(44, 8, 1, 3500),
(44, 2, 1, 3000),
(44, 2, 1, 3000),
(45, 11, 1, 1000),
(45, 11, 1, 1000),
(45, 8, 1, 3500),
(45, 2, 1, 3000),
(45, 8, 1, 3500),
(46, 8, 1, 3500),
(46, 8, 1, 3500),
(46, 2, 1, 3000),
(46, 2, 1, 3000),
(47, 20, 1, 1500),
(47, 8, 1, 3500),
(47, 8, 1, 3500),
(48, 12, 1, 2000),
(48, 8, 1, 3500),
(48, 8, 1, 3500),
(48, 2, 1, 3000),
(48, 2, 1, 3000),
(49, 7, 1, 2800),
(49, 12, 1, 2000),
(50, 8, 1, 3500),
(54, 8, 1, 3500),
(54, 8, 1, 3500),
(54, 8, 1, 3500),
(64, 20, 1, 1500),
(68, 20, 1, 1500),
(71, 2, 1, 3000),
(71, 2, 1, 3000),
(71, 2, 1, 3000),
(71, 2, 1, 3000),
(71, 2, 1, 3000),
(71, 2, 1, 3000),
(71, 2, 1, 3000),
(71, 2, 1, 3000),
(71, 2, 1, 3000),
(71, 2, 1, 3000),
(72, 21, 1, 2000),
(72, 21, 1, 2000),
(72, 21, 1, 2000),
(72, 21, 1, 2000),
(72, 21, 1, 2000),
(74, 20, 1, 1500),
(74, 20, 1, 1500),
(74, 20, 1, 1500),
(74, 20, 1, 1500),
(74, 20, 1, 1500),
(74, 20, 1, 1500),
(74, 20, 1, 1500),
(74, 20, 1, 1500),
(74, 20, 1, 1500),
(74, 20, 1, 1500),
(77, 11, 1, 1000),
(77, 11, 1, 1000),
(77, 11, 1, 1000),
(78, 12, 1, 2000),
(81, 20, 1, 1500),
(81, 20, 1, 1500),
(81, 20, 1, 1500),
(81, 20, 1, 1500),
(84, 12, 1, 2000),
(84, 12, 1, 2000),
(85, 11, 1, 1000),
(85, 11, 1, 1000),
(85, 11, 1, 1000),
(88, 11, 1, 1000),
(88, 11, 1, 1000),
(88, 11, 1, 1000),
(91, 23, 1, 8000),
(91, 23, 1, 8000),
(95, 12, 1, 2000),
(96, 8, 1, 3500),
(97, 20, 1, 1500),
(98, 21, 1, 2000),
(100, 22, 1, 3500),
(102, 21, 1, 2000),
(102, 21, 1, 2000),
(103, 12, 1, 2000),
(103, 12, 1, 2000),
(103, 12, 1, 2000),
(104, 12, 1, 2000),
(105, 9, 1, 1800),
(106, 21, 1, 2000),
(106, 21, 1, 2000),
(106, 21, 1, 2000),
(107, 25, 1, 1000),
(107, 25, 1, 1000),
(107, 25, 1, 1000),
(107, 25, 1, 1000),
(107, 25, 1, 1000),
(108, 25, 1, 1000),
(108, 25, 1, 1000),
(108, 25, 1, 1000),
(109, 25, 1, 1000),
(109, 25, 1, 1000),
(109, 25, 1, 1000),
(109, 25, 1, 1000),
(109, 25, 1, 1000),
(109, 21, 1, 2000),
(110, 21, 1, 2000),
(110, 21, 1, 2000),
(110, 21, 1, 2000),
(110, 21, 1, 2000),
(110, 21, 1, 2000),
(111, 25, 1, 1000),
(111, 25, 1, 1000),
(111, 25, 1, 1000),
(111, 25, 1, 1000),
(111, 25, 1, 1000),
(112, 21, 1, 2000),
(112, 21, 1, 2000),
(112, 21, 1, 2000),
(113, 21, 1, 2000),
(114, 21, 1, 2000),
(115, 22, 1, 3500),
(116, 25, 1, 1000),
(116, 25, 1, 1000),
(116, 25, 1, 1000),
(116, 25, 1, 1000),
(116, 25, 1, 1000),
(116, 25, 1, 1000),
(116, 25, 1, 1000),
(116, 25, 1, 1000),
(116, 25, 1, 1000),
(116, 25, 1, 1000),
(114, 25, 1, 1000),
(114, 25, 1, 1000),
(114, 25, 1, 1000),
(114, 25, 1, 1000),
(114, 25, 1, 1000),
(114, 25, 1, 1000),
(114, 25, 1, 1000),
(114, 25, 1, 1000),
(114, 25, 1, 1000),
(114, 25, 1, 1000),
(114, 25, 1, 1000),
(114, 25, 1, 1000),
(117, 21, 1, 2000),
(117, 20, 1, 1500),
(118, 25, 1, 1000),
(118, 25, 1, 1000),
(118, 25, 1, 1000),
(118, 25, 1, 1000),
(118, 8, 1, 3500),
(118, 8, 1, 3500),
(119, 25, 1, 1000),
(119, 25, 1, 1000),
(119, 25, 1, 1000),
(119, 25, 1, 1000),
(119, 25, 1, 1000),
(119, 25, 1, 1000),
(120, 25, 1, 1000),
(120, 25, 1, 1000),
(120, 25, 1, 1000),
(120, 25, 1, 1000),
(120, 25, 1, 1000);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `pedidos_modificadores`
--

CREATE TABLE `pedidos_modificadores` (
  `id` int(11) NOT NULL,
  `id_ped` int(11) NOT NULL,
  `id_pro` int(11) NOT NULL,
  `id_mod` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `pedidos_modificadores`
--

INSERT INTO `pedidos_modificadores` (`id`, `id_ped`, `id_pro`, `id_mod`) VALUES
(1, 42, 8, 1),
(2, 43, 8, 1),
(3, 44, 8, 1),
(4, 44, 2, 3),
(5, 45, 8, 1),
(6, 45, 2, 3),
(7, 46, 8, 1),
(8, 46, 2, 3),
(9, 47, 8, 2),
(10, 48, 8, 1),
(11, 48, 2, 3),
(15, 71, 2, 3),
(16, 71, 2, 3),
(17, 71, 2, 3),
(18, 71, 2, 3),
(19, 71, 2, 3),
(20, 71, 2, 3),
(21, 71, 2, 3),
(22, 71, 2, 3),
(23, 71, 2, 3),
(24, 71, 2, 3),
(26, 91, 23, 7),
(27, 91, 23, 9),
(29, 91, 23, 7),
(30, 91, 23, 8),
(31, 96, 8, 2),
(32, 100, 22, 6),
(35, 118, 8, 1),
(37, 118, 8, 1);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `pedidos_online`
--

CREATE TABLE `pedidos_online` (
  `id_online` int(11) NOT NULL,
  `fecha` datetime NOT NULL DEFAULT current_timestamp(),
  `cliente` varchar(50) NOT NULL,
  `total` int(11) NOT NULL,
  `estado` varchar(20) DEFAULT 'pendiente',
  `tipo_pedido` varchar(20) DEFAULT 'retiro',
  `notas` text DEFAULT NULL,
  `pedido_json` text DEFAULT NULL,
  `turno_id` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `pedidos_online`
--

INSERT INTO `pedidos_online` (`id_online`, `fecha`, `cliente`, `total`, `estado`, `tipo_pedido`, `notas`, `pedido_json`, `turno_id`) VALUES
(1, '2026-03-22 18:43:59', 'cxcx', 20000, 'listo', 'retiro', NULL, NULL, 0),
(2, '2026-03-22 18:45:41', 'kjl', 2800, 'listo', 'retiro', NULL, NULL, 0),
(3, '2026-03-22 19:15:22', 'sdf', 5600, 'listo', '', NULL, NULL, 0),
(4, '2026-03-22 19:16:24', 'kjl', 1500, 'listo', '', NULL, NULL, 0),
(5, '2026-03-22 19:19:07', 'sdf', 8500, 'listo', '', NULL, NULL, 0),
(6, '2026-03-22 19:29:36', 'erer', 5400, 'listo', 'Comer en local', NULL, NULL, 0),
(7, '2026-03-22 20:17:25', 'cxcx', 9900, 'listo', 'Retiro', NULL, NULL, 0),
(8, '2026-03-23 10:16:06', 'fghfghfgh', 14000, 'listo', 'Comer en local', NULL, NULL, 0),
(9, '2026-03-23 10:22:01', 'dfgd', 4800, 'listo', 'Retiro', NULL, NULL, 0),
(10, '2026-03-23 11:11:06', 'jjjhjh', 3300, 'listo', 'Retiro', NULL, NULL, 0),
(11, '2026-03-23 14:54:05', 'dfsdf', 5600, 'listo', 'Retiro', NULL, NULL, 0),
(12, '2026-03-23 15:15:11', 'adfaf', 5600, 'listo', 'Retiro', NULL, NULL, 0),
(13, '2026-03-24 14:09:28', 'adfaf', 8200, 'listo', 'Retiro', NULL, NULL, 0),
(14, '2026-03-24 14:21:58', 'dfsdf', 12000, 'listo', 'Comer en local', NULL, NULL, 0),
(15, '2026-03-24 19:02:24', 'sdf', 4200, 'listo', 'Retiro', NULL, NULL, 0),
(16, '2026-03-24 19:36:20', 'erer', 6000, 'listo', 'Retiro', NULL, NULL, 0),
(17, '2026-03-24 19:37:20', 'adfaf', 8000, 'listo', 'Retiro', NULL, NULL, 0),
(18, '2026-03-24 19:44:49', 'cxcx', 5000, 'listo', 'Comer en local', NULL, NULL, 0),
(19, '2026-03-24 19:51:33', 'cxcx', 6500, 'listo', 'Comer en local', NULL, NULL, 0),
(20, '2026-03-25 10:35:36', 'dfsdf', 3500, 'listo', 'Comer en local', NULL, NULL, 0),
(21, '2026-03-25 11:10:45', 'sdf', 3000, 'listo', 'Retiro', NULL, NULL, 0),
(22, '2026-03-25 11:37:19', 'cxcx', 1000, 'listo', 'Retiro', NULL, NULL, 0),
(23, '2026-03-25 11:49:47', 'dfsdf', 5000, 'listo', 'Retiro', NULL, NULL, 0),
(24, '2026-03-25 11:51:45', 'erer', 1000, 'listo', 'Retiro', NULL, NULL, 0),
(25, '2026-03-25 11:54:08', 'kjl', 5000, 'listo', 'Comer en local', NULL, NULL, 0),
(26, '2026-03-25 11:55:43', 'cxcx', 5000, 'listo', 'Retiro', NULL, NULL, 0),
(27, '2026-03-25 11:56:41', 'kjl', 5000, 'listo', 'Retiro', NULL, NULL, 0),
(28, '2026-03-25 11:57:36', 'paula', 4000, 'listo', 'Retiro', NULL, NULL, 0),
(29, '2026-03-25 12:12:18', 'dfgdgfdfg', 5000, 'listo', 'Retiro', 'muchas servilletas', NULL, 0),
(30, '2026-03-25 12:20:58', 'paula', 5000, 'pendiente', 'Retiro', '[Coca cola: con hielo] \r\nmuchas servilletas', NULL, 0),
(31, '2026-03-25 12:24:27', 'paula', 1000, 'pendiente', 'Retiro', '[Coca cola: con hielo] \r\nmuhcas servilletas', NULL, 0),
(32, '2026-03-25 12:32:23', 'paula', 5000, 'pendiente', 'Retiro', 'sin servilletas', NULL, 0),
(33, '2026-03-25 12:36:55', 'paula', 4000, 'pendiente', 'Retiro', '• Coca cola: Sin hielo\r\nsin servilletas', NULL, 0),
(34, '2026-03-25 12:57:16', 'Paula', 4000, 'pendiente', 'Retiro', '• Coca cola: sin hielo\r\ncon servilletas', '[{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[{\"id\":\"0\",\"nombre\":\"sin hielo\"}]},{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[{\"id\":\"0\",\"nombre\":\"sin hielo\"}]},{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[{\"id\":\"0\",\"nombre\":\"sin hielo\"}]},{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[{\"id\":\"0\",\"nombre\":\"sin hielo\"}]}]', 0),
(35, '2026-03-25 13:03:50', 'erer', 3000, 'listo', 'Retiro', '• Coca cola: sin hielo\r\nnada', '[{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[{\"id\":\"0\",\"nombre\":\"sin hielo\"}]},{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[{\"id\":\"0\",\"nombre\":\"sin hielo\"}]},{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[{\"id\":\"0\",\"nombre\":\"sin hielo\"}]}]', 0),
(36, '2026-03-25 15:03:53', 'dfsdf', 4000, 'pendiente', 'Retiro', 'dssfdfsdfsdfsdfsdfsfd', '[{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[]},{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[]},{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[]},{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[]}]', 0),
(37, '2026-03-25 18:08:18', 'cxcx', 5000, 'listo', 'Retiro', '• Coca cola: cccccccccc\r\nasasass', '[{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[{\"id\":\"0\",\"nombre\":\"cccccccccc\"}]},{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[{\"id\":\"0\",\"nombre\":\"cccccccccc\"}]},{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[{\"id\":\"0\",\"nombre\":\"cccccccccc\"}]},{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[{\"id\":\"0\",\"nombre\":\"cccccccccc\"}]},{\"idPro\":25,\"nombre\":\"Coca cola\",\"mods\":[{\"id\":\"0\",\"nombre\":\"cccccccccc\"}]}]', 11);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `pedidos_online_detalle`
--

CREATE TABLE `pedidos_online_detalle` (
  `id_online` int(11) NOT NULL,
  `id_pro` int(11) NOT NULL,
  `cantidad` int(11) NOT NULL,
  `precio` int(11) NOT NULL,
  `notas_producto` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `pedidos_online_detalle`
--

INSERT INTO `pedidos_online_detalle` (`id_online`, `id_pro`, `cantidad`, `precio`, `notas_producto`) VALUES
(1, 10, 1, 1500, NULL),
(1, 10, 1, 1500, NULL),
(1, 9, 1, 1800, NULL),
(1, 9, 1, 1800, NULL),
(4, 20, 1, 1500, NULL),
(6, 12, 1, 2000, NULL),
(6, 12, 1, 2000, NULL),
(7, 8, 1, 3500, NULL),
(8, 11, 1, 1000, NULL),
(8, 23, 1, 8000, NULL),
(9, 21, 1, 2000, NULL),
(9, 7, 1, 2800, NULL),
(10, 9, 1, 1800, NULL),
(10, 10, 1, 1500, NULL),
(13, 11, 1, 1000, NULL),
(13, 11, 1, 1000, NULL),
(13, 12, 1, 2000, NULL),
(16, 11, 1, 1000, NULL),
(16, 11, 1, 1000, NULL),
(16, 11, 1, 1000, NULL),
(16, 11, 1, 1000, NULL),
(16, 11, 1, 1000, NULL),
(16, 11, 1, 1000, NULL),
(17, 21, 1, 2000, NULL),
(17, 21, 1, 2000, NULL),
(17, 21, 1, 2000, NULL),
(17, 21, 1, 2000, NULL),
(18, 11, 5, 5000, NULL),
(19, 8, 1, 3500, NULL),
(19, 2, 1, 3000, NULL),
(20, 8, 1, 3500, NULL),
(21, 25, 1, 1000, NULL),
(21, 11, 2, 2000, NULL),
(22, 25, 1, 1000, NULL),
(23, 25, 5, 5000, NULL),
(24, 25, 1, 1000, NULL),
(25, 25, 5, 5000, NULL),
(26, 25, 5, 5000, NULL),
(27, 25, 5, 5000, NULL),
(28, 25, 4, 4000, NULL),
(29, 25, 5, 5000, NULL),
(30, 25, 5, 5000, NULL),
(31, 25, 1, 1000, NULL),
(32, 25, 5, 5000, NULL),
(33, 25, 4, 4000, NULL),
(34, 25, 4, 4000, 'sin hielo'),
(35, 25, 3, 3000, 'sin hielo'),
(36, 25, 4, 4000, ''),
(37, 25, 5, 5000, 'cccccccccc');

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `pedidos_online_modificadores`
--

CREATE TABLE `pedidos_online_modificadores` (
  `id_online` int(11) NOT NULL,
  `id_pro` int(11) NOT NULL,
  `id_mod` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `pedidos_online_modificadores`
--

INSERT INTO `pedidos_online_modificadores` (`id_online`, `id_pro`, `id_mod`) VALUES
(7, 8, 1),
(7, 8, 2),
(8, 23, 7),
(8, 23, 8),
(8, 23, 11),
(19, 8, 1),
(19, 2, 0),
(21, 25, 0),
(21, 11, 0),
(22, 25, 0),
(23, 25, 0),
(24, 25, 0),
(25, 25, 0),
(26, 25, 0),
(27, 25, 0),
(28, 25, 0),
(29, 25, 0),
(30, 25, 0),
(31, 25, 0),
(32, 25, 0),
(33, 25, 0);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `productos`
--

CREATE TABLE `productos` (
  `id_pro` int(4) NOT NULL,
  `nombre` varchar(50) NOT NULL,
  `precio` int(4) NOT NULL,
  `id_cat` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `productos`
--

INSERT INTO `productos` (`id_pro`, `nombre`, `precio`, `id_cat`) VALUES
(2, 'italiano', 3000, 2),
(7, 'Sandwich Jamón', 2800, 2),
(8, 'Churrasco', 3500, 2),
(9, 'Papas Fritas', 1800, 5),
(10, 'Papas Medianas', 1500, 5),
(11, 'Helados', 1000, 6),
(12, 'Salame', 2000, 7),
(20, 'chocolate', 1500, 10),
(21, 'pino', 2000, 11),
(22, 'Hamburguesa Italiana', 3500, 12),
(23, 'Promo 1', 8000, 13),
(25, 'Coca cola', 1000, 15);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `recetas`
--

CREATE TABLE `recetas` (
  `id` int(11) NOT NULL,
  `id_pro` int(11) NOT NULL,
  `id_ing` int(11) NOT NULL,
  `cantidad` int(11) NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `turnos`
--

CREATE TABLE `turnos` (
  `id_turno` int(11) NOT NULL,
  `nombre` varchar(50) NOT NULL,
  `inicio` datetime NOT NULL,
  `fin` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `turnos`
--

INSERT INTO `turnos` (`id_turno`, `nombre`, `inicio`, `fin`) VALUES
(1, 'Turno', '2026-03-21 15:24:49', '2026-03-21 15:24:56'),
(2, 'Turno', '2026-03-21 15:25:19', '2026-03-21 15:25:24'),
(3, '21 marzo', '2026-03-21 15:26:31', '2026-03-21 15:34:03'),
(4, '21 marzo', '2026-03-21 15:34:12', '2026-03-21 15:46:11'),
(5, '21/03', '2026-03-21 15:46:19', '2026-03-21 16:12:14'),
(6, '21 marzo', '2026-03-21 16:12:18', '2026-03-21 16:25:41'),
(7, 'Turno', '2026-03-21 16:25:43', '2026-03-21 16:25:48'),
(8, '21 marzo', '2026-03-21 16:25:51', '2026-03-21 16:35:40'),
(9, 'Turno', '2026-03-21 16:35:41', '2026-03-21 16:35:47'),
(10, '21 marzo', '2026-03-21 16:35:50', '2026-03-25 10:36:36'),
(11, 'Turno', '2026-03-25 10:36:40', NULL);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `usuarios`
--

CREATE TABLE `usuarios` (
  `id_usr` int(4) NOT NULL,
  `nombre` varchar(30) NOT NULL,
  `password` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Volcado de datos para la tabla `usuarios`
--

INSERT INTO `usuarios` (`id_usr`, `nombre`, `password`) VALUES
(1, 'Hola', '$2a$10$PpRmmcNr5pEQzmMxmOS0cevCobFy0LHrSgleMaJTrkbWe35BvgKsu');

--
-- Índices para tablas volcadas
--

--
-- Indices de la tabla `categorias`
--
ALTER TABLE `categorias`
  ADD PRIMARY KEY (`id_cat`);

--
-- Indices de la tabla `configuracion`
--
ALTER TABLE `configuracion`
  ADD PRIMARY KEY (`clave`);

--
-- Indices de la tabla `ingredientes`
--
ALTER TABLE `ingredientes`
  ADD PRIMARY KEY (`id_ing`);

--
-- Indices de la tabla `mesas`
--
ALTER TABLE `mesas`
  ADD PRIMARY KEY (`id_mesa`);

--
-- Indices de la tabla `modificadores`
--
ALTER TABLE `modificadores`
  ADD PRIMARY KEY (`id_mod`),
  ADD KEY `id_pro` (`id_pro`);

--
-- Indices de la tabla `pedidos`
--
ALTER TABLE `pedidos`
  ADD PRIMARY KEY (`id_ped`);

--
-- Indices de la tabla `pedidos_detalle`
--
ALTER TABLE `pedidos_detalle`
  ADD KEY `id_pro` (`id_pro`),
  ADD KEY `pedidos_detalle_ibfk_1` (`id_ped`);

--
-- Indices de la tabla `pedidos_modificadores`
--
ALTER TABLE `pedidos_modificadores`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id_ped` (`id_ped`),
  ADD KEY `id_mod` (`id_mod`);

--
-- Indices de la tabla `pedidos_online`
--
ALTER TABLE `pedidos_online`
  ADD PRIMARY KEY (`id_online`);

--
-- Indices de la tabla `pedidos_online_detalle`
--
ALTER TABLE `pedidos_online_detalle`
  ADD KEY `id_online` (`id_online`),
  ADD KEY `id_pro` (`id_pro`);

--
-- Indices de la tabla `pedidos_online_modificadores`
--
ALTER TABLE `pedidos_online_modificadores`
  ADD KEY `id_online` (`id_online`);

--
-- Indices de la tabla `productos`
--
ALTER TABLE `productos`
  ADD PRIMARY KEY (`id_pro`),
  ADD KEY `id_cat` (`id_cat`);

--
-- Indices de la tabla `recetas`
--
ALTER TABLE `recetas`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id_pro` (`id_pro`),
  ADD KEY `id_ing` (`id_ing`);

--
-- Indices de la tabla `turnos`
--
ALTER TABLE `turnos`
  ADD PRIMARY KEY (`id_turno`);

--
-- Indices de la tabla `usuarios`
--
ALTER TABLE `usuarios`
  ADD PRIMARY KEY (`id_usr`);

--
-- AUTO_INCREMENT de las tablas volcadas
--

--
-- AUTO_INCREMENT de la tabla `categorias`
--
ALTER TABLE `categorias`
  MODIFY `id_cat` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16;

--
-- AUTO_INCREMENT de la tabla `ingredientes`
--
ALTER TABLE `ingredientes`
  MODIFY `id_ing` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=15;

--
-- AUTO_INCREMENT de la tabla `mesas`
--
ALTER TABLE `mesas`
  MODIFY `id_mesa` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=12;

--
-- AUTO_INCREMENT de la tabla `modificadores`
--
ALTER TABLE `modificadores`
  MODIFY `id_mod` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=14;

--
-- AUTO_INCREMENT de la tabla `pedidos`
--
ALTER TABLE `pedidos`
  MODIFY `id_ped` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=121;

--
-- AUTO_INCREMENT de la tabla `pedidos_modificadores`
--
ALTER TABLE `pedidos_modificadores`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=39;

--
-- AUTO_INCREMENT de la tabla `pedidos_online`
--
ALTER TABLE `pedidos_online`
  MODIFY `id_online` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=38;

--
-- AUTO_INCREMENT de la tabla `productos`
--
ALTER TABLE `productos`
  MODIFY `id_pro` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=26;

--
-- AUTO_INCREMENT de la tabla `recetas`
--
ALTER TABLE `recetas`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT de la tabla `turnos`
--
ALTER TABLE `turnos`
  MODIFY `id_turno` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=12;

--
-- AUTO_INCREMENT de la tabla `usuarios`
--
ALTER TABLE `usuarios`
  MODIFY `id_usr` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- Restricciones para tablas volcadas
--

--
-- Filtros para la tabla `modificadores`
--
ALTER TABLE `modificadores`
  ADD CONSTRAINT `modificadores_ibfk_1` FOREIGN KEY (`id_pro`) REFERENCES `productos` (`id_pro`);

--
-- Filtros para la tabla `pedidos_detalle`
--
ALTER TABLE `pedidos_detalle`
  ADD CONSTRAINT `pedidos_detalle_ibfk_1` FOREIGN KEY (`id_ped`) REFERENCES `pedidos` (`id_ped`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `pedidos_detalle_ibfk_2` FOREIGN KEY (`id_pro`) REFERENCES `productos` (`id_pro`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `pedidos_modificadores`
--
ALTER TABLE `pedidos_modificadores`
  ADD CONSTRAINT `pedidos_modificadores_ibfk_1` FOREIGN KEY (`id_ped`) REFERENCES `pedidos` (`id_ped`),
  ADD CONSTRAINT `pedidos_modificadores_ibfk_2` FOREIGN KEY (`id_mod`) REFERENCES `modificadores` (`id_mod`);

--
-- Filtros para la tabla `pedidos_online_detalle`
--
ALTER TABLE `pedidos_online_detalle`
  ADD CONSTRAINT `pedidos_online_detalle_ibfk_1` FOREIGN KEY (`id_online`) REFERENCES `pedidos_online` (`id_online`),
  ADD CONSTRAINT `pedidos_online_detalle_ibfk_2` FOREIGN KEY (`id_pro`) REFERENCES `productos` (`id_pro`);

--
-- Filtros para la tabla `pedidos_online_modificadores`
--
ALTER TABLE `pedidos_online_modificadores`
  ADD CONSTRAINT `pedidos_online_modificadores_ibfk_1` FOREIGN KEY (`id_online`) REFERENCES `pedidos_online` (`id_online`);

--
-- Filtros para la tabla `productos`
--
ALTER TABLE `productos`
  ADD CONSTRAINT `productos_ibfk_1` FOREIGN KEY (`id_cat`) REFERENCES `categorias` (`id_cat`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `recetas`
--
ALTER TABLE `recetas`
  ADD CONSTRAINT `recetas_ibfk_1` FOREIGN KEY (`id_pro`) REFERENCES `productos` (`id_pro`),
  ADD CONSTRAINT `recetas_ibfk_2` FOREIGN KEY (`id_ing`) REFERENCES `ingredientes` (`id_ing`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
