<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddSentEmailIdToComplaints extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('complaints', function (Blueprint $table) { 
            $table->integer('sent_email_id')->unsigned()->after('type');
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
        Schema::table('complaints', function (Blueprint $table) {
            $table->dropForeign('sent_email_id');
            $table->dropColumn('sent_email_id');
        });
    }
}
