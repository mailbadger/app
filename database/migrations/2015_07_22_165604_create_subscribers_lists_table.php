<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateSubscribersListsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('subscribers_lists', function (Blueprint $table) {
            $table->increments('id');
            $table->integer('subscriber_id')->unsigned();
            $table->foreign('subscriber_id')->references('id')->on('subscribers');
            $table->integer('list_id')->unsigned();
            $table->foreign('list_id')->references('id')->on('lists');
            $table->unique(['subscriber_id', 'list_id']);
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::drop('subscribers_lists');
    }
}
