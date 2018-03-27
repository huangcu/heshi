<template>
  <div>
    <h3>合适代理
      <span class="accountlevelstars">{{ fromAgentLevel }}</span>
    </h3>
    <p class="inpage-navi-box">
      <button class="inpage-navi active" id="inpage-navi_generalinfo" v-on:click="goToSection('generalinfo')">基本信息</button>
      <span id="new-iconbox">
        <img id="newicon-mysite" src="/_images/constant/new-icon.png" />
      </span>
      <button class="inpage-navi mywebsitebtn" id="inpage-navi_mywebsite" v-on:click="goToSection('mywebsite')">我的网站</button>
      <button class="inpage-navi" id="inpage-navi_mypoints" v-on:click="goToSection('mypoints')">我的积分</button>
      <button class="inpage-navi" id="inpage-navi_myclients" v-on:click="goToSection('myclients')">我的客人</button>
      <button class="inpage-navi" id="inpage-navi_clientorders" v-on:click="goToSection('clientorders')">客人的预定</button>
      <button class="inpage-navi" id="inpage-navi_coupon" v-on:click="goToSection('coupon')">打折／代金券</button>
      <button class="inpage-navi" id="inpage-navi_personalinfo" v-on:click="goToSection('personalinfo')">个人信息</button>
      <a class="inpage-navi" id="listdownloadbtn" href="stocklist-csv.php">下载货单 &raquo;</a>
    </p>
    <!-- offerticket box -->
    <div class="generalinfo discountbox" id="generalinfo">
      <h3 class="g-info-title">您好，{{ name }}
        <form method="post" action="" id="logoutform">
          <input type="hidden" name="logout" value="confirmed" />
          <input type="submit" value="登出账号" name="logoutbtn" id="logoutbtn" />
        </form>
      </h3>

      <div v-if="agent_level_locked=='NO'">
        <h4 class="g-info-content">您目前的代理级别: {{ agentLevel}}</h4>
        <h4 class="g-info-content">购买折扣: {{ priceforagent(agentLevel, 1) }}</h4>
        <h4>我的积分：(一年内) 总购买金额: <strong>{{ total_point }}</strong>欧元；钻石购买总颗数: <strong style="display:inline-block; margin-right:15px;">{{ total_dia_num }}</strong>
          <button id="mypointdetailbtn" href="/myaccount#section-mypoints" v-on:click="goToSection('mypoints')">详情 &raquo;</button>
        </h4>
        <p class="ticket-question" v-on:click="showRegulation()">折扣、积分、返点规则 &raquo;</p>
        <div class="ticket-answer">
          合适代理折扣：
          <ul>
            <li>一级代理（享有 9 折价格）: 一年购买12颗钻石以下 或 总购买金额在4万欧元以下</li>
            <li>二级代理（享有8.5折价格）: 一年内购买超过12颗钻石而且总购买金额超过4万欧元</li>
            <li>三级代理（享有8.3折价格）: 一年内购买超过20颗钻石而且总购买金额超过10万欧元</li>
          </ul>
          代理积分、返点规则:
          <ul>
            <li>
            一级代理如果一年购买12颗钻石以上而且总购买金额超过4万欧元，将获得一年内总购买金额的5%的返点，并且升为二级代理
            </li>
            <li>
            二级代理如果一年购买超过20颗钻石，将获得一年内（如果上次返点在一年内，则从上次返点开始到现在）总购买金额的2%的返点，并且升为三级代理。如果一年内购买低于12颗或总金额低于4万欧元，则降为1级代理。
            </li>
            <li>
            三级级代理已享有最高折扣，不再返点。如果一年内购买低于20颗钻石或总金额低于10万欧元，相应的降为二级或一级代理。
            </li>
          </ul>
        </div>
      </div>
      <div v-else>
        <h4>您目前的代理级别: {{ agentLevel }} 级 &nbsp; &nbsp; &nbsp; 购买折扣为: {{ priceforagent(agentLevel, 1) }}</h4>
      </div>

      <div id="myreferenceticketbox">
        <p id="titleformyreference" v-on:click="openRecoContent()">
        轻松推广 增加客人 &raquo;
        </p>
        <div id="reco-contentbox">
          由您推荐的注册用户都是您的客人。您的客人仍然可以向下继续推荐客人。<br />
          您的客人推荐的客人如果购买合适商品，您同样享有<strong>1％</strong>的返点。您客人推荐的客人再推荐客人，您依然享有1%的返点。以此类推，惊喜无限！<br />
          <div class="share-ticket-tools">
            <span class="glyphicon glyphicon-hand-right"></span> 推荐方法:
            <p id="toolexpl">点击打开下面链接中的页面，然后把该页面分享到微信朋友圈(或发送给朋友)。也可以拷贝这个链接，邮件等方式发给朋友。就这么简单！</p>
            <p><a target="_blank" :href="'http://beyoudiamond.com/login?ref='+account+'X'+offer_ticket" id="recommenduserurl">http://beyoudiamond.com/login.php?ref={{accountID}}'X'{{offer_ticket}}</a></p>
          </div>
        </div>
      </div><!-- end roco-contentbox -->
    </div>  <!-- end offerticket box -->

    <!-- my website -->
    <div class="mywebsite" id="mywebsite">
      <!-- $filename='reseller'.$accountID.'/wp-content/themes/flatsome-child/account-configuration.php'; -->
      <div v-if="filenameexist">
        <h3>您的销售利器，我们已经为您准备好！</h3>
        <p><br />合适总部可以专门为代理配置属于代理自己的网站。</p>
        <p>网站中不出现「合适钻石」的品牌信息，而使用代理自己的品牌。</p>
        <p>钻石和首饰的库存和合适钻石保持一致，显示的价格可以由代理自己调整。</p>
        <p>开通网站和详情请咨询合适总部<strong>何丽丽</strong><br /><br /></p>
      </div>
      <div v-else>
        <div class="mysitemanagingbloc">
          <div v-if="feedbackmessageofrequest">
            <h3>{{ feedbackmessageofrequest }}</h3>
            <p class="sitestatuswords">
              <span class="glyphicon glyphicon-thumbs-up"></span>您的网站已经开通 ({{ domaintoshow }})
              <a target="_blank" :href="active_siteurl">查看
                <span class="glyphicon glyphicon-new-window"></span>
              </a>
            </p>
            <form action="" method="post">
              <p>
                <label>商品价格：</label>
                <input type="number" step="0.01" name="resellerpriceratio" id="resellerpriceratio" :value="resellerpriceratio" /> x 合适原始价格
              </p>
              <input type="submit" value="保存" />
            </form>
            <div class="sitemanageformbox">
              <form name="loginform" id="loginform" :action="active_siteurl+'/wp-login.php'" method="post" target="_blank">
                <p style="display:none;">
                  <label for="user_pass">UserName<br>
                  <input type="text" name="log" id="user_login" class="input" value="manager" size="20"></label>
                </p>
                <p style="display:none;">
                  <label for="user_pass">Password<br>
                  <input type="password" name="pwd" id="user_pass" class="input" value="BeyouR!2017" size="20"></label>
                </p>
                <p class="forgetmenot" style="display:none;">
                  <label for="rememberme"><input name="rememberme" type="checkbox" id="rememberme" value="forever" checked="checked"> Remember Me</label>
                </p>
                <p class="submit">
                  <input type="submit" name="wp-submit" id="wp-submit" class="button button-primary button-large" value="管理网站内容">
                  <input type="hidden" name="redirect_to" :value="active_siteurl">
                  <input type="hidden" name="testcookie" value="1">
                </p>
              </form>
              <p style="font-size:12px"><span class="glyphicon glyphicon-info-sign"></span> 如果需要登录，用户名: manager 密码: BeyouR!2017</p>
            </div>
          </div>
          <div v-else>
            <div v-if="userdefineddomainname!==''">
              <p class="sitestatuswords">
                <span class="glyphicon glyphicon-ok"></span> 我们已经收到您提交的信息，请稍等，我们会在24小时内为您开通。</p>
              <p>
                <label>网址(域名)选择：</label>
                <input type="text" name="sitedomain" disabled="disabled" id="sitedomain" :value="userdefineddomainname" />
              </p>
              <p>
                <label>品牌名称：</label>
                <input type="text" name="sitetitle" id="sitetitle" disabled="disabled" value="site_title" />
              </p>
            </div>
            <div v-else>
              <p class="sitestatuswords">请填写您的网站信息并提交，填好后我们会在24小时内为您开通。</p>
              <!-- start of apply site -->
              <form action="" method="post">
                <p>
                  <label>网址(域名)选择：</label>
                  <input type="text" name="sitedomain" id="sitedomain" value="" />
                  <select name="domainsuffix" id="domainsuffix">
                    <option value="com">.com</option>
                    <option value="be">.be</option>
                    <option value="info">.info</option>
                    <option value="net">.net</option>
                    <option value="org">.org</option>
                    <option value="online">.online</option>
                    <option value="biz">.biz</option>
                    <option value="me">.me</option>
                  </select>
                  <button type="button" id="domaincheckbtn" v-on:click="checkdomainavailability()">查询是否可用
                    <span class="glyphicon glyphicon-search"></span>
                  </button>
                  <span class="note">
                    <span class="glyphicon glyphicon-exclamation-sign"></span> 注：提交后不可更改
                  </span>
                </p>
                <span v-if="r_domain_active" class="domainstatuswords">域名:{{ r_website_domain }} 申请中，请稍等
                </span>
                <span v-else-if="r_domain_active=='OCCU'">域名:{{ r_website_domain }} 已经被占用，请另选一个
                </span>
                <span v-else>域名:{{ r_website_domain }}
                </span>
                <p>
                  <label>品牌名称：</label>
                  <input type="text" name="sitetitle" id="sitetitle" :value="site_title" />
                  <span class="note">
                    <span class="glyphicon glyphicon-exclamation-sign"></span> 注：提交后不可更改
                  </span>
                </p>
                <p>
                  <label>商品价格：</label>
                  <input type="number" step="0.01" name="resellerpriceratio" id="resellerpriceratio" :value="resellerpriceratio" /> x 合适原始价格
                  <span class="note">
                    <span class="glyphicon glyphicon-info-sign"></span> 注：提交后可以随时更改
                  </span>
                </p>
                <input type="submit" value="提交" />
              </form><!-- end of apply site -->
            </div>
          </div>
        </div>

        <div class="tech_support_bloc">
          <h4>技术支持</h4>
          <h5>如何更改网站上的内容? </h5>
          <p>操作非常简单，请看一下我们为您录制的两段几分钟的视频，您即能轻松掌握。</p>
          <p>
            <a target="_blank" href="http://beyoudiamond.com/press/%E4%BB%A3%E8%B4%AD%E7%BD%91%E7%AB%99%E4%BD%BF%E7%94%A8%E4%BB%8B%E7%BB%8D-%E9%A6%96%E9%A1%B5-2/">请点击这里，</a>查看视频解说 (上集)</p>
          <p>
            <a target="_blank" href="http://beyoudiamond.com/press/%E4%BB%A3%E8%B4%AD%E7%BD%91%E7%AB%99%E4%BD%BF%E7%94%A8%E4%BB%8B%E7%BB%8D-%E5%85%B6%E4%BB%96%E8%AF%B4%E6%98%8E-small/">请点击这里，</a>查看视频解说 (下集)</p>
          <p>更多的问题，请点击右侧的“在线咨询”联系我们，我们会给您详尽的解答</p>
        </div>
        <br class="clear" />
      </div>
    </div>
    <!-- orders -->
    <div class="orders" id="clientorders">
      <h3>客人的订单</h3>
      <div v-if="transactionfound">
        <p>目前尚无订单可以显示</p>
      </div>
      <div v-else>
        <ul>
          <!-- $thisTransactionDiaList.='<span class="glyphicon glyphicon-ok"></span>';
            $thisTransactionDiaList.='<span class="glyphicon glyphicon-hourglass"></span>';
                  $thisTransactionJewList.='<span class="glyphicon glyphicon-ok"></span>';
          $thisTransactionJewList.='<span class="glyphicon glyphicon-hourglass"></span>';-->
          <li class="transactionpiece" :class="crr_tr_heshi_status" :id="'transactionpiece_'+crr_transaction_id">
            <span class="buyername">
              <label><span class="glyphicon glyphicon-user"></span> 预定人</label>
              <a target="_blank" :href="'userhistory.php?id='+crr_tr_buyer_id">{{ crr_buyer_name }} &raquo; </a>
            </span>
            <span class="transactiondiamonds">
              <label>预定的钻石(ID)</label>[{{thisTransactionDiaList }}</span>
            <span class="transactionjewelry">
              <label>预定的首饰(ID)</label>[{{thisTransactionJewList }}</span>
            <span class="transactiontotalprice">
              <label>总价格</label>{{ crr_tr_total_price }} €</span>
            <span class="transactionheshibi">
              <label>使用合适币</label>[{{crr_tr_heshibi_amount }}</span>
            <span class="transactionfinalprice">
              <label>最终应付款</label>{{ crr_tr_final_price }} €</span>
            <span class="agentprice">
              <label>代理利润</label>{{ crr_tr_agent_profit }} €</span>
            <span class="transactiontime">下单时间: [{{crr_tr_created }}</span>
            <p class="transactionstatusactionsbox">
              <div v-if="crr_tr_heshi_status=='PROCESSING'">
                <span class="precessinginfo">
                  <span class="glyphicon glyphicon-th-list"></span> 合适总部正在处理订单
                </span>
              </div>
              <div v-else-if="crr_tr_heshi_status=='COMPLETE'">
                <div v-if="crr_tr_agent_status=='PROCESSING'">
                  <button class="transactioncompletebtn" :id="'transactioncompletebtn_'+crr_transaction_id" v-on:click="completeTransaction(crr_transaction_id)">
                    <span class="glyphicon glyphicon-check"></span> 标记为交易完成
                  </button>
                </div>
                <div v-else>
                  <span class="tr_complete_words">
                    <span class="glyphicon glyphicon-ok"></span> 交易已完成
                  </span>
                  <button class="archivebtn" :id="'transactionarchive_'+crr_transaction_id" v-on:click="archiveTransaction(crr_transaction_id)">
                    <span class="glyphicon glyphicon-inbox"></span> 存档
                  </button>
                </div>
              </div>
              <!-- no end tag for p -->
          </li>
        </ul>
      </div>
    </div>

    <!-- history -->
    <div class="history" id="mypoints">
      <h3>购买历史纪录</h3>
      <div class="historybox history-diamond">
        <h4>钻石购买历史纪录</h4>
        <ul v-if="found">
          <!-- require_once('_includes/functions/detail_fordiamond_byshape.php');
          $buyernametxt='<p class="orderredbytxt"><span class="thetxtofusername"><span class="glyphicon glyphicon-user"></span> 预定人: '.$nameofuser.'</span></p>'; -->
          <li v-for="(item, index) in items" :key="index" :class="'dia-order-history'+listclass" :id="'diaorder-'+diamondid"> {{ buyernametxt }}
            <span class="dia-description">{{ detail_forDiamond_byShape($crr_dia_shape, 'NAMECN') }}钻石: {{crr_dia_weight }} 克拉，{{ crr_dia_color}} 色， 净度 {{ crr_dia_clarity }}</span>
            <span v-if="jewelry_name!=''" class="jewelry-description">配饰: {{jewelry_name}}</span>
            <span v-else>无</span>
            <span class="totalpricecrr">总价格: {{ sold_price_total}} >欧元</span>
            <span class="orderedtime">
              <span v-if="order_at<last_year">超一年 不积分</span>
              <span v-else class="glyphicon glyphicon-ok">积分有效</span>
            </span>
            <br class="clear" />
          </li>
        </ul>
        <ul v-else>
          <li class="no-order">没有可积分的钻石购买纪录</li>
        </ul>
        <ul v-if="nopointfound">
          <!-- require_once('_includes/functions/detail_fordiamond_byshape.php'); -->
          <!-- $buyernametxt='<p class="orderredbytxt"><span class="thetxtofusername"><span class="glyphicon glyphicon-user"></span> 预定人: '.$nameofuser.'</span></p>'; -->
          <li v-for="(item, index) in items" :key="index" id="no-profit-title-box">&darr; 不能积分的交易纪录 &darr; {{ buyernametxt }}
            <span class="dia-description">{{ detail_forDiamond_byShape($crr_dia_shape, 'NAMECN') }}钻石: {{crr_dia_weight }} 克拉，{{ crr_dia_color}} 色， 净度 {{ crr_dia_clarity }}</span>
            <span v-if="jewelry_name!=''" class="jewelry-description">配饰: {{jewelry_name}}</span>
            <span v-else>无</span>
            <span class="totalpricecrr">总价格: {{ sold_price_total}} >欧元</span>
            <span class="orderedtime">该交易不可用于积分 {{ ordered_at }}</span>
            <br class="clear" />
          </li>
        </ul>
      </div>

      <div class="historybox history-diamond">
        <h4>首饰购买历史记录</h4>
        <ul v-if="found">
          <!-- $buyernametxt='<p class="orderredbytxt"><span class="thetxtofusername"><span class="glyphicon glyphicon-user"></span> 预定人: '.$nameofuser.'</span></p>'; -->
          <li v-for="(item, index) in items" :key="index" :class="'dia-order-history'+listclass"> {{ buyernametxt }}
            <span class="jewelry-name">{{ crr_jew_category }}</span>
            <span class="jewelry-description">{{ crr_jew_name }}</span>
            <span class="totalpricecrr">总价格: {{ sold_price_total}} 欧元</span>
            <span class="orderedtime">
              <span v-if="order_at<last_year">超一年 不积分</span>
              <span v-else class="glyphicon glyphicon-ok">积分有效(空托不累计数量)</span>
            </span>
            <br class="clear" />
          </li>
        </ul>
        <ul v-else>
          <li class="no-order">尚未购买任何首饰</li>
        </ul>
      </div>

      <div class="historysummary">
        <p>您一年内总的购买金额为 <span class="thepricevalue"> {{ total_price_one_year }} </span>欧元</p>
        <p>您一年内共购买<span class="thepricevalue">{{ total_piece_one_year }}</span>件合适商品。</p>
        <p v-if="agent_level_locked=='NO'" class="agentlevelwords"> {{ agent_level_words }}</p>
        <p>您的全部购买金额为 <span class="thepricevalue"> {{ total_price }}</span>欧元</p>
      </div>
    </div>

    <!-- history recommended users-->
    <div class="history-recommenedusers" id="myclients">
      <h3>我的客户</h3>
      <div class="theuserslist">
        <ul v-if="found" class="recommenduserslist">
          <li v-for="(item, index) in items" :key="index">
            <span class="clientpic" :style="'background-image:url('+ru_icon+')'"></span>
            <span class="ru_name">{{ crr_ru_name }}</span>
            <span  v-if="totalordernum!=''" class="neworder"> ( <span class="glyphicon glyphicon-gift"></span> 有新预订 ) </span>
            <span class="usertotalamout">一年内消费总额:  {{ crr_user_total_all }} 欧元 (共购买 {{ crr_user_total_pieces }}件商品)
              <a target="_blank" class="userhistorybtn" :href="'userhistory.php?id='+ru_id">
                <span class="userdetailbtn">详情 &raquo;</span>
              </a>
            </span>
            <a class="chatlinker" :href="'/message-agent.php?c='+ru_id">
              <span class="glyphicon glyphicon-comment"></span> 消息
              <span v-if="usermn!=''" class='messagemn newmessage' :class="'messagemn'+mn_class_add"> [{{ usermn }} 则新消息] </span>
              <span v-else class="'messagemn"> [无新消息] </span>
            </a>
            <br class="clear" />
          </li>
        </ul>
        <ul v-else class="recommenduserslist">
          <p>您尚未推荐任何用户</p>
        </ul>
      </div>
    </div>

    <!-- history extended users-->
    <div class="history-recommenedusers extendedclients">
      <h3>客户发展的客人</h3>
      <div class="theuserslist">
        <ul v-if="found" class="recommenduserslist">
          <li v-for="(item, index) in items" :key="index">
            <span class="ru_name">{{ crr_rue_name }}</span>
            <span v-if="crr_usere_total_all!=''" class="ru_price newpoints">(总消费: {{crr_usere_total_all }}欧元)</span>
            <a target="_blank" class="userhistorybtn" :href="'userhistory.php?id='+rue_id">详情 &raquo;</a>
            <br class="clear" />
          </li>
        </ul>
        <ul v-else class="recommenduserslist">
          <p>您的客人尚未推荐任何用户</p>
        </ul>
      </div>
    </div>

    <!-- discount ticket -->
    <div class="generalinfo heshibi-box discountticketbox" id="coupon">
      <h3>我的打折券</h3>
      <div v-if="dtfount">
        <span v-for="(item, index) in items" :key="index" class="discountticketcontainer">
          <img :src="'/_images/constant/discount'+crr_dt_value+'.jpg'" class="dt-pic" />
          <span class="dt-time">
            <span class="glyphicon glyphicon-time"></span> 有效期:
            <!-- date('Y年m月d日', strtotime(date('Y-m-d', strtotime($crr_dt_created)+8*60*60).' +1 Month'));  -->
          </span>
        </span>
      </div>
      <div v-else>
        <span> 没有可以使用的打折券</span>
      </div>
      <p class="discountticketexplwords">
        您购买商品时，系统会自动选择最早的一张打折卡添加到您的订单中。<br />
        打折券的折扣和您自己折扣不能累加，以两者中折扣较高的为准。
      </p>
      <h4 class="heshibitranferbtnbox dt">
        <a href="senddiscountticket.php">
          <span class="glyphicon glyphicon-hand-right"></span> 发打折券给我的客人 &raquo;
        </a>
      </h4>
      <p class="heshibitransfer-expl">
        <span class="glyphicon glyphicon-info-sign"></span> 您可以给您的客人发打折券作为他们购买时的折扣凭证
      </p>
    </div>

    <!-- heshibi -->
    <div class="generalinfo heshibi-box">
      <h3>我的合适币（代金劵）</h3>
      <div class="my-heshibi">
        <div class="my-heshibi-amount-box">
          余额: <span id="heshibi-balance-value">{{ balance }}</span>€
        </div>
      </div>
      <p>合适币余额: {{ balance }} 欧元, {{ EuroToYuan( balance) }}元人民币, {{ EuroToDollar(balance) }}美元</p>
      <h4 class="heshibitranferbtnbox">
        <a href="heshibitransfer.php">
          <span class="glyphicon glyphicon-hand-right"></span> 发合适币给我的客人 &raquo;
        </a>
      </h4>
      <p class="heshibitransfer-expl">
        <span class="glyphicon glyphicon-info-sign"></span> 如果客人一星期内不使用，金额自动返回
      </p>
    </div>
  </div>
</template>
