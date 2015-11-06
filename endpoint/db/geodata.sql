CREATE TABLE IF NOT EXISTS `locations` (
  `locid` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `latitude` float(6,4) NOT NULL,
  `longitude` float(7,4) NOT NULL,
  `ipcount` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`locid`),
  KEY `latlon` (`latitude`,`longitude`,`ipcount`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;
