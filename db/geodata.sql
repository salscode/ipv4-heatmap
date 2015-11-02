CREATE TABLE IF NOT EXISTS `locations` (
  `locid` mediumint(8) unsigned NOT NULL,
  `latitude` float(8,4) NOT NULL,
  `longitude` float(8,4) NOT NULL,
  `ipcount` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`locid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
