<template>
  <div>
    <vue-title title="添加管理人员"></vue-title>
    <h1>添加普通管理人员</h1>
    <p style="text-align:center;">普通管理人员权限：仅网站库存</p>
    <p v-for="(error, index) in registerErrors" :key="index" class='alert'>{{ error }}</p>
    <!-- <p v-if="error" class='alert'>{{ error}}</p>
    <p v-if="message_db" class='alert'>{{ message_db}}</p> -->
    <div id="contentwrapper">
      <form @submit.prevent="newAdminUser()" id="addnewadmin">
        <p>
          <label>姓名: </label><span>(必填)</span>
          <input name="name" id="name" type="text" v-model="realName" v-validate="'required'" class="formbox" size="200"/>
        </p>
        <p>
          <label>电话号码: </label><span>(必填)</span>
          <input name="cellphone" id="cellphone" type="number" v-model.number="cellphone" class="formbox" size="200"/>
        </p>
        <p>
          <label>用户名: </label><span>(英文字母，必填)</span>
          <input name="username" id="username" type="text" v-model="username" v-validate="'required'" class="formbox" size="200"/>
        </p>
        <p>
          <label>密码: </label><span>(英文字母或数字，必填)</span>
          <input name="password" id="password" type="password" v-model="password" v-validate="'required'" class="formbox" size="200"/>
        </p>
        <p>
          <label>确认密码: </label><span>(英文字母或数字，必填)</span>
          <input name="password" id="password2" type="password" v-model="password2" v-validate="'required'" class="formbox" size="200"/>
        </p>
        <p>
          <label>管理员级别: </label><span>(数字1-10必填)</span>
          <input name="adminlevel" id="adminlevel" type="number" v-model.number="adminLevel" class="formbox" size="200"/>
        </p>
        <p>
          <label>头像:</label>
          <input type="file" multiple id="image" class="image_selecting" name="icon" @change.prevent="validateIconFile($event.target.files)" accept="image/*" />
          <button @click.prevent = "clearIconFile">clear</button>
        </p>
        <input name="newadminformsubmit" type="submit" id="newadminbtn" value="提交" />
      </form>
      <!-- <div class="picture" style="position:relative;">
        <div class="picture-preview-box">
        </div>
        <form action="uploadingimage.php" method="post" enctype="multipart/form-data" id="uploadingadminicon" class="uploadImage">
          <input type="hidden" name="picturewhere" value="adminicon" />
          <label>添加头像:</label>
          <input type="file" name="image" id="image" class="image_selecting">
        </form>
        <br style="clear:both;" />
        <div class="progress" style="width:300px; position:relative; top:-320px; left:488px;">
          <p>上传中</p>
          <div class="bar"></div >
          <div class="percent">0%</div >
        </div>
      </div>
      <p style="padding-top:30px; text-align:center; padding-left:0; background:none;">
        <button onclick="submittheform()" type="button" id="submit_article" style="font-size:18px; font-weight:bold; padding:12px 50px; background-color:#393; color:#FFF; border-width:1px;"><span class="glyphicon glyphicon-send" style="color:#FFF; left:0;"></span> &nbsp;  提交 </button>
      </p> -->
    </div>
    <div id="status" style="display:none;"></div>
    <div id="iframecontainer" style="display:none;"></div>
    <br /><br /><br /><br /><br /><br />
  </div>
</template>
<script src="./NewAdmin.js"></script>
<style src="./NewAdmin.css" scoped></style>
