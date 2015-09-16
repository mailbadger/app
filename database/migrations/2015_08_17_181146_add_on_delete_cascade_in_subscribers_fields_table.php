<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddOnDeleteCascadeInSubscribersFieldsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('subscribers_fields', function (Blueprint $table) {
            $table->dropForeign('subscribers_fields_field_id_foreign');
            $table->foreign('field_id')->references('id')->on('fields')->onDelete('cascade');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('subscribers_fields', function (Blueprint $table) {
            $table->dropForeign('subscribers_fields_field_id_foreign');
            $table->foreign('field_id')->references('id')->on('fields');
        });
    }
}
