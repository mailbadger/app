<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateBouncesTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('bounces', function (Blueprint $table) {
            $table->increments('id');
            $table->string('recipient');
            $table->string('sender');
            $table->string('type');
            $table->string('sub_type');
            $table->string('action');
            $table->dateTime('timestamp');
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
        Schema::drop('bounces');
    }
}
