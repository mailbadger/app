<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddSentEmailIdToBounces extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('bounces', function (Blueprint $table) {
            $table->integer('sent_email_id')->unsigned()->after('action');
            $table->foreign('sent_email_id')->references('id')->on('sent_emails');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('bounces', function (Blueprint $table) {
            $table->dropForeign('sent_email_id');
            $table->dropColumn('sent_email_id');
        });
    }
}
